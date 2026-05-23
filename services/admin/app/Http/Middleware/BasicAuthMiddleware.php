<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class BasicAuthMiddleware
{
    public function handle(Request $request, Closure $next): Response
    {
        $expectedUser = env('ADMIN_USER');
        $expectedPassword = env('ADMIN_PASSWORD');

        // 未設定時は fail closed
        if (empty($expectedUser) || empty($expectedPassword)) {
            return $this->unauthorized();
        }

        $user = $request->getUser();
        $password = $request->getPassword();

        if (
            !hash_equals($expectedUser, (string) $user) ||
            !hash_equals($expectedPassword, (string) $password)
        ) {
            return $this->unauthorized();
        }

        return $next($request);
    }

    private function unauthorized(): Response
    {
        return response('Unauthorized', 401, [
            'WWW-Authenticate' => 'Basic realm="Lang War Admin"',
        ]);
    }
}
