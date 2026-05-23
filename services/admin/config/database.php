<?php

use Illuminate\Support\Str;

return [
    'default' => 'pgsql',

    'connections' => [
        'pgsql' => (function () {
            $base = [
                'driver'  => 'pgsql',
                'charset' => 'utf8',
                'prefix'  => '',
                'schema'  => 'public',
                'sslmode' => env('DB_SSLMODE', 'prefer'),
            ];

            if ($url = env('DATABASE_URL')) {
                $parsed = parse_url($url);
                return array_merge($base, [
                    'host'     => $parsed['host'] ?? '127.0.0.1',
                    'port'     => $parsed['port'] ?? '5432',
                    'database' => ltrim($parsed['path'] ?? '', '/'),
                    'username' => $parsed['user'] ?? '',
                    'password' => isset($parsed['pass']) ? urldecode($parsed['pass']) : '',
                ]);
            }

            return array_merge($base, [
                'host'     => env('DB_HOST', '127.0.0.1'),
                'port'     => env('DB_PORT', '5432'),
                'database' => env('DB_DATABASE', 'lang_war'),
                'username' => env('DB_USERNAME', 'lang_war'),
                'password' => env('DB_PASSWORD', ''),
            ]);
        })(),
    ],

    'migrations' => [
        'table' => 'migrations',
        'update_date_on_publish' => true,
    ],
];
