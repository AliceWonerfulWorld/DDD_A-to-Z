<?php

namespace App\Http\Controllers;

use App\Models\AdminLog;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class LogController
{
    public function index(): Response
    {
        $operationLogs = AdminLog::orderByDesc('created_at')->limit(200)->get();

        $appLogLines = $this->readAppLog(200);

        return Inertia::render('Logs/Index', [
            'operationLogs' => $operationLogs,
            'appLogLines'   => $appLogLines,
        ]);
    }

    private function readAppLog(int $lines): array
    {
        $path = storage_path('logs/laravel.log');

        if (!file_exists($path)) {
            return [];
        }

        $file = new \SplFileObject($path, 'r');
        $file->seek(PHP_INT_MAX);
        $total = $file->key();

        $start  = max(0, $total - $lines);
        $result = [];
        $file->seek($start);

        while (!$file->eof()) {
            $line = $file->current();
            if ($line !== '' && $line !== false) {
                $result[] = rtrim((string) $line);
            }
            $file->next();
        }

        return $result;
    }
}
