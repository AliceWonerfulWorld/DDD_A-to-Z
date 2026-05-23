<?php

namespace Tests\Feature;

use Tests\TestCase;

class BasicAuthTest extends TestCase
{
    public function test_it_rejects_requests_without_basic_auth_credentials(): void
    {
        $this->get('/guilds')
            ->assertUnauthorized()
            ->assertHeader('WWW-Authenticate', 'Basic realm="Lang War Admin"');
    }

    public function test_it_allows_requests_with_valid_basic_auth_credentials(): void
    {
        $this->withServerVariables($this->basicAuthHeaders())
            ->get('/guilds')
            ->assertOk();
    }
}
