<?php

namespace Tests\Feature;

use App\Models\Guild;
use Tests\TestCase;

class GuildControllerTest extends TestCase
{
    public function test_store_creates_a_guild_and_records_an_admin_log(): void
    {
        $this->withServerVariables($this->basicAuthHeaders())
            ->post('/guilds', [
                'slug' => 'backend-guild',
                'name' => 'Backend Guild',
                'description' => 'API and infrastructure specialists',
                'icon' => 'server',
                'color' => '#1A2B3C',
                'sort_order' => 2,
            ])
            ->assertRedirect('/guilds')
            ->assertSessionHas('success', 'ギルドを作成しました');

        $guild = Guild::where('slug', 'backend-guild')->first();

        $this->assertNotNull($guild);
        $this->assertNotEmpty($guild->id);
        $this->assertSame('Backend Guild', $guild->name);
        $this->assertSame(0, $guild->current_exp);
        $this->assertSame(1, $guild->guild_level);

        $this->assertDatabaseHas('admin_logs', [
            'action' => 'created',
            'target_type' => 'guild',
            'target_id' => 'backend-guild',
        ]);
    }

    public function test_store_rejects_duplicate_guild_slugs(): void
    {
        Guild::create([
            'id' => '01HV0000000000000000000000',
            'slug' => 'backend-guild',
            'name' => 'Backend Guild',
            'description' => 'Existing guild',
            'icon' => 'server',
            'color' => '#1A2B3C',
            'sort_order' => 1,
            'created_at' => now(),
            'updated_at' => now(),
        ]);

        $this->withServerVariables($this->basicAuthHeaders())
            ->from('/guilds/create')
            ->post('/guilds', [
                'slug' => 'backend-guild',
                'name' => 'Backend Guild 2',
                'description' => 'Duplicate guild',
                'icon' => 'server',
                'color' => '#1A2B3C',
                'sort_order' => 2,
            ])
            ->assertRedirect('/guilds/create')
            ->assertSessionHasErrors('slug');
    }
}
