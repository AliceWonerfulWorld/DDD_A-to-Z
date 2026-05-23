import core.thread : Thread;
import core.time : Duration, MonoTime, dur;
import std.algorithm.comparison : max, min;
import std.algorithm.sorting : sort;
import std.algorithm.searching : endsWith, startsWith;
import std.array : appender;
import std.conv : to;
import std.datetime.stopwatch : StopWatch;
import std.exception : enforce;
import std.format : formattedWrite;
import std.math : ceil;
import std.net.curl : CurlException, HTTP, HTTPStatusException;
import std.stdio : stderr, writeln;
import std.string : toLower;

struct Config
{
    string target = "http://localhost:8080/healthz";
    Duration duration = dur!"seconds"(5);
    uint concurrency = 2;
    Duration timeout = dur!"seconds"(2);
    string output = "json";
    bool allowExternal = false;
    bool help = false;
}

struct RequestSample
{
    ushort statusCode;
    long latencyMs;
    bool transportOk;
    string error;
}

struct LoadTestReport
{
    string target;
    long durationMs;
    uint concurrency;
    ulong totalRequests;
    ulong successRequests;
    ulong failedRequests;
    ulong[string] statusCounts;
    ulong[string] errorCounts;
    double averageLatencyMs;
    long p95LatencyMs;
    long p99LatencyMs;
}

int main(string[] args)
{
    try
    {
        auto config = parseArgs(args);
        if (config.help)
        {
            writeln(usage());
            return 0;
        }

        validateConfig(config);
        auto report = runLoadTest(config);

        if (config.output == "json")
        {
            writeln(toJson(report));
        }
        else
        {
            writeln(toText(report));
        }

        return 0;
    }
    catch (Exception e)
    {
        stderr.writeln("error: ", e.msg);
        stderr.writeln();
        stderr.writeln(usage());
        return 1;
    }
}

Config parseArgs(string[] args)
{
    Config config;

    for (size_t i = 1; i < args.length; i++)
    {
        auto arg = args[i];

        if (arg == "--")
        {
            continue;
        }

        string readValue(string flag)
        {
            auto prefix = flag ~ "=";
            if (arg.startsWith(prefix))
            {
                return arg[prefix.length .. $];
            }

            enforce(i + 1 < args.length, flag ~ " requires a value");
            i++;
            return args[i];
        }

        if (arg == "--help" || arg == "-h")
        {
            config.help = true;
        }
        else if (arg == "--allow-external")
        {
            config.allowExternal = true;
        }
        else if (arg == "--target" || arg.startsWith("--target="))
        {
            config.target = readValue("--target");
        }
        else if (arg == "--duration" || arg.startsWith("--duration="))
        {
            config.duration = parseDuration(readValue("--duration"));
        }
        else if (arg == "--concurrency" || arg.startsWith("--concurrency="))
        {
            config.concurrency = readValue("--concurrency").to!uint;
        }
        else if (arg == "--timeout" || arg.startsWith("--timeout="))
        {
            config.timeout = parseDuration(readValue("--timeout"));
        }
        else if (arg == "--output" || arg.startsWith("--output="))
        {
            config.output = readValue("--output").toLower;
        }
        else
        {
            throw new Exception("unknown option: " ~ arg);
        }
    }

    return config;
}

Duration parseDuration(string value)
{
    enforce(value.length > 0, "duration must not be empty");

    string numberPart;
    string unit;

    if (value.endsWith("ms"))
    {
        numberPart = value[0 .. $ - 2];
        unit = "ms";
    }
    else if (value.endsWith("s") || value.endsWith("m") || value.endsWith("h"))
    {
        numberPart = value[0 .. $ - 1];
        unit = value[$ - 1 .. $];
    }
    else
    {
        throw new Exception("duration must use one of: ms, s, m, h");
    }

    auto amount = numberPart.to!long;
    enforce(amount > 0, "duration must be greater than zero");

    if (unit == "ms")
    {
        return dur!"msecs"(amount);
    }
    if (unit == "s")
    {
        return dur!"seconds"(amount);
    }
    if (unit == "m")
    {
        return dur!"minutes"(amount);
    }
    if (unit == "h")
    {
        return dur!"hours"(amount);
    }

    assert(0, "unreachable duration unit");
}

void validateConfig(Config config)
{
    enforce(config.target.startsWith("http://") || config.target.startsWith("https://"),
        "--target must start with http:// or https://");
    enforce(config.concurrency >= 1, "--concurrency must be at least 1");
    enforce(config.concurrency <= 200, "--concurrency is capped at 200");
    enforce(config.duration <= dur!"hours"(1), "--duration is capped at 1h");
    enforce(config.timeout <= dur!"minutes"(5), "--timeout is capped at 5m");
    enforce(config.output == "json" || config.output == "text",
        "--output must be json or text");

    if (!config.allowExternal)
    {
        enforce(isLocalTarget(config.target),
            "external targets require --allow-external");
    }
}

bool isLocalTarget(string target)
{
    auto host = hostFromUrl(target);
    return host == "localhost" || host == "127.0.0.1" || host == "::1" || host == "[::1]";
}

string hostFromUrl(string target)
{
    auto lowerTarget = target.toLower;
    size_t start;

    if (lowerTarget.startsWith("http://"))
    {
        start = "http://".length;
    }
    else if (lowerTarget.startsWith("https://"))
    {
        start = "https://".length;
    }
    else
    {
        return "";
    }

    auto rest = target[start .. $];

    if (rest.length > 0 && rest[0] == '[')
    {
        foreach (idx, ch; rest)
        {
            if (ch == ']')
            {
                return rest[0 .. idx + 1].toLower;
            }
        }

        return rest.toLower;
    }

    size_t end = rest.length;

    foreach (idx, ch; rest)
    {
        if (ch == '/' || ch == ':' || ch == '?' || ch == '#')
        {
            end = idx;
            break;
        }
    }

    return rest[0 .. end].toLower;
}

LoadTestReport runLoadTest(Config config)
{
    auto deadline = MonoTime.currTime + config.duration;
    auto samplesByWorker = new RequestSample[][](config.concurrency);
    auto threads = new Thread[](config.concurrency);

    foreach (i; 0 .. config.concurrency)
    {
        immutable workerIndex = i;
        threads[workerIndex] = new Thread(delegate void() {
            samplesByWorker[workerIndex] = runWorker(config, deadline);
        });
        threads[workerIndex].start();
    }

    foreach (thread; threads)
    {
        thread.join();
    }

    RequestSample[] samples;
    foreach (workerSamples; samplesByWorker)
    {
        samples ~= workerSamples;
    }

    return summarize(config, samples);
}

RequestSample[] runWorker(Config config, MonoTime deadline)
{
    RequestSample[] samples;

    while (MonoTime.currTime < deadline)
    {
        samples ~= performRequest(config.target, config.timeout);
    }

    return samples;
}

RequestSample performRequest(string target, Duration timeout)
{
    StopWatch watch;
    watch.start();

    try
    {
        auto http = HTTP(target);
        http.method = HTTP.Method.get;
        http.connectTimeout = timeout;
        http.dataTimeout = timeout;
        http.operationTimeout = timeout;
        http.maxRedirects = 0;
        http.onReceive = (ubyte[] data) { return data.length; };
        http.perform();
        watch.stop();

        return RequestSample(http.statusLine.code, watch.peek.total!"msecs", true, "");
    }
    catch (HTTPStatusException e)
    {
        watch.stop();
        return RequestSample(cast(ushort) e.status, watch.peek.total!"msecs", true, "");
    }
    catch (CurlException e)
    {
        watch.stop();
        return RequestSample(0, watch.peek.total!"msecs", false, e.msg);
    }
    catch (Exception e)
    {
        watch.stop();
        return RequestSample(0, watch.peek.total!"msecs", false, e.msg);
    }
}

LoadTestReport summarize(Config config, RequestSample[] samples)
{
    LoadTestReport report;
    report.target = config.target;
    report.durationMs = config.duration.total!"msecs";
    report.concurrency = config.concurrency;
    report.totalRequests = samples.length;

    long[] latencies;
    long latencyTotal;

    foreach (sample; samples)
    {
        latencies ~= sample.latencyMs;
        latencyTotal += sample.latencyMs;

        if (sample.transportOk)
        {
            auto key = sample.statusCode.to!string;
            report.statusCounts[key] = report.statusCounts.get(key, 0) + 1;

            if (sample.statusCode >= 200 && sample.statusCode < 400)
            {
                report.successRequests++;
            }
            else
            {
                report.failedRequests++;
            }
        }
        else
        {
            report.failedRequests++;
            auto key = sample.error.length == 0 ? "transport_error" : sample.error;
            report.errorCounts[key] = report.errorCounts.get(key, 0) + 1;
        }
    }

    if (samples.length > 0)
    {
        report.averageLatencyMs = cast(double) latencyTotal / cast(double) samples.length;
        sort(latencies);
        report.p95LatencyMs = percentile(latencies, 95);
        report.p99LatencyMs = percentile(latencies, 99);
    }

    return report;
}

long percentile(long[] sortedValues, int percentileRank)
{
    if (sortedValues.length == 0)
    {
        return 0;
    }

    auto rawIndex = cast(long) ceil((cast(double) percentileRank / 100.0)
        * cast(double) sortedValues.length) - 1;
    auto index = cast(size_t) min(max(rawIndex, 0), cast(long) sortedValues.length - 1);
    return sortedValues[index];
}

string toJson(LoadTestReport report)
{
    auto outp = appender!string;
    outp.put("{");
    formattedWrite(outp, `"target":"%s",`, jsonEscape(report.target));
    formattedWrite(outp, `"durationMs":%s,`, report.durationMs);
    formattedWrite(outp, `"concurrency":%s,`, report.concurrency);
    formattedWrite(outp, `"totalRequests":%s,`, report.totalRequests);
    formattedWrite(outp, `"successRequests":%s,`, report.successRequests);
    formattedWrite(outp, `"failedRequests":%s,`, report.failedRequests);
    formattedWrite(outp, `"statusCounts":%s,`, mapToJson(report.statusCounts));
    formattedWrite(outp, `"errorCounts":%s,`, mapToJson(report.errorCounts));
    formattedWrite(outp, `"latencyMs":{"avg":%.2f,"p95":%s,"p99":%s}`,
        report.averageLatencyMs, report.p95LatencyMs, report.p99LatencyMs);
    outp.put("}");
    return outp.data;
}

string mapToJson(ulong[string] values)
{
    auto outp = appender!string;
    auto keys = values.keys;
    sort(keys);

    outp.put("{");
    foreach (idx, key; keys)
    {
        if (idx > 0)
        {
            outp.put(",");
        }
        formattedWrite(outp, `"%s":%s`, jsonEscape(key), values[key]);
    }
    outp.put("}");
    return outp.data;
}

string jsonEscape(string value)
{
    auto outp = appender!string;

    foreach (ch; value)
    {
        if (ch == '"')
        {
            outp.put(`\"`);
        }
        else if (ch == '\\')
        {
            outp.put(`\\`);
        }
        else if (ch == '\n')
        {
            outp.put(`\n`);
        }
        else if (ch == '\r')
        {
            outp.put(`\r`);
        }
        else if (ch == '\t')
        {
            outp.put(`\t`);
        }
        else if (ch < 0x20)
        {
            formattedWrite(outp, `\u%04x`, cast(uint) ch);
        }
        else
        {
            outp.put(ch);
        }
    }

    return outp.data;
}

string toText(LoadTestReport report)
{
    auto outp = appender!string;
    formattedWrite(outp, "target: %s\n", report.target);
    formattedWrite(outp, "durationMs: %s\n", report.durationMs);
    formattedWrite(outp, "concurrency: %s\n", report.concurrency);
    formattedWrite(outp, "totalRequests: %s\n", report.totalRequests);
    formattedWrite(outp, "successRequests: %s\n", report.successRequests);
    formattedWrite(outp, "failedRequests: %s\n", report.failedRequests);
    formattedWrite(outp, "statusCounts: %s\n", mapToJson(report.statusCounts));
    formattedWrite(outp, "errorCounts: %s\n", mapToJson(report.errorCounts));
    formattedWrite(outp, "latencyMs.avg: %.2f\n", report.averageLatencyMs);
    formattedWrite(outp, "latencyMs.p95: %s\n", report.p95LatencyMs);
    formattedWrite(outp, "latencyMs.p99: %s", report.p99LatencyMs);
    return outp.data;
}

string usage()
{
    return q"USAGE
Lang War Dlang defensive load tester

Usage:
  pnpm --filter @lang-war/loadtest-d dev -- --target http://localhost:8080/healthz --duration 10s --concurrency 5 --timeout 2s --output json

Options:
  --target <url>         Target URL. Defaults to http://localhost:8080/healthz.
  --duration <duration>  Test duration. Supports ms, s, m, h. Defaults to 5s.
  --concurrency <n>      Number of worker threads. Defaults to 2. Capped at 200.
  --timeout <duration>   Per-request timeout. Defaults to 2s.
  --output <json|text>   Output format. Defaults to json.
  --allow-external       Required for non-localhost targets.
  --help                 Show this help.

Success criteria:
  HTTP 200-399 responses are counted as successful requests.
  Redirects are not followed, so 3xx responses are reported as-is and counted as success.
USAGE";
}

unittest
{
    assert(parseDuration("250ms") == dur!"msecs"(250));
    assert(parseDuration("2s") == dur!"seconds"(2));
    assert(parseDuration("3m") == dur!"minutes"(3));
    assert(parseDuration("1h") == dur!"hours"(1));
}

unittest
{
    assert(isLocalTarget("http://localhost:8080/healthz"));
    assert(isLocalTarget("https://127.0.0.1/test"));
    assert(isLocalTarget("http://[::1]:8080/healthz"));
    assert(!isLocalTarget("https://example.com/healthz"));
    assert(!isLocalTarget("http://localhost.example.com/healthz"));
}

unittest
{
    assert(percentile([1, 2, 3, 4, 5], 95) == 5);
    assert(percentile([1, 2, 3, 4, 5], 50) == 3);
    assert(percentile([], 99) == 0);
}

unittest
{
    auto samples = [
        RequestSample(200, 10, true, ""),
        RequestSample(500, 30, true, ""),
        RequestSample(0, 20, false, "timeout")
    ];
    Config config;
    auto report = summarize(config, samples);
    assert(report.totalRequests == 3);
    assert(report.successRequests == 1);
    assert(report.failedRequests == 2);
    assert(report.statusCounts["200"] == 1);
    assert(report.statusCounts["500"] == 1);
    assert(report.errorCounts["timeout"] == 1);
    assert(report.p95LatencyMs == 30);
}
