<?php

namespace Tests;

use Illuminate\Foundation\Testing\TestCase as BaseTestCase;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;

abstract class TestCase extends BaseTestCase
{
    protected function setUp(): void
    {
        parent::setUp();

        $this->configureInMemoryDatabase();
        $this->rebuildSchema();
    }

    protected function basicAuthHeaders(string $user = 'admin', string $password = 'secret'): array
    {
        return [
            'PHP_AUTH_USER' => $user,
            'PHP_AUTH_PW' => $password,
        ];
    }

    private function configureInMemoryDatabase(): void
    {
        config()->set('database.default', 'testing');
        config()->set('database.connections.testing', [
            'driver' => 'sqlite',
            'database' => ':memory:',
            'prefix' => '',
            'foreign_key_constraints' => true,
        ]);

        DB::purge('testing');
        DB::reconnect('testing');
    }

    private function rebuildSchema(): void
    {
        Schema::dropIfExists('admin_logs');
        Schema::dropIfExists('guilds');
        Schema::dropIfExists('point_types');

        Schema::create('guilds', function ($table): void {
            $table->text('id')->primary();
            $table->text('slug')->unique();
            $table->text('name');
            $table->text('description');
            $table->text('icon');
            $table->text('color');
            $table->integer('sort_order');
            $table->bigInteger('current_exp')->default(0);
            $table->integer('guild_level')->default(1);
            $table->timestampTz('created_at');
            $table->timestampTz('updated_at');
        });

        Schema::create('point_types', function ($table): void {
            $table->text('code');
            $table->text('language')->default('');
            $table->text('label');
            $table->primary(['code', 'language']);
        });

        Schema::create('admin_logs', function ($table): void {
            $table->id();
            $table->text('action');
            $table->text('target_type');
            $table->text('target_id');
            $table->json('payload');
            $table->timestampTz('created_at');
        });
    }
}
