<?php

namespace Tests\Feature;

use App\Models\PointType;
use Tests\TestCase;

class PointTypeControllerTest extends TestCase
{
    public function test_store_creates_a_point_type_and_records_an_admin_log(): void
    {
        $this->withServerVariables($this->basicAuthHeaders())
            ->post('/point-types', [
                'code' => 'SP',
                'language' => 'PHP',
                'label' => 'PHP Skill Point',
            ])
            ->assertRedirect('/point-types')
            ->assertSessionHas('success', 'ポイントタイプを作成しました');

        $this->assertDatabaseHas('point_types', [
            'code' => 'SP',
            'language' => 'PHP',
            'label' => 'PHP Skill Point',
        ]);
        $this->assertDatabaseHas('admin_logs', [
            'action' => 'created',
            'target_type' => 'point_type',
            'target_id' => 'SP:PHP',
        ]);
    }

    public function test_store_rejects_duplicate_code_and_language_pairs(): void
    {
        PointType::create([
            'code' => 'SP',
            'language' => 'PHP',
            'label' => 'PHP Skill Point',
        ]);

        $this->withServerVariables($this->basicAuthHeaders())
            ->from('/point-types/create')
            ->post('/point-types', [
                'code' => 'SP',
                'language' => 'PHP',
                'label' => 'Duplicate PHP Skill Point',
            ])
            ->assertRedirect('/point-types/create')
            ->assertSessionHasErrors('code');
    }

    public function test_update_changes_only_the_matching_code_and_language_pair(): void
    {
        PointType::create([
            'code' => 'SP',
            'language' => 'PHP',
            'label' => 'PHP Skill Point',
        ]);
        PointType::create([
            'code' => 'SP',
            'language' => 'Go',
            'label' => 'Go Skill Point',
        ]);

        $this->withServerVariables($this->basicAuthHeaders())
            ->put('/point-types', [
                'code' => 'SP',
                'language' => 'PHP',
                'label' => 'PHP Mastery Point',
            ])
            ->assertRedirect('/point-types')
            ->assertSessionHas('success', 'ポイントタイプを更新しました');

        $this->assertDatabaseHas('point_types', [
            'code' => 'SP',
            'language' => 'PHP',
            'label' => 'PHP Mastery Point',
        ]);
        $this->assertDatabaseHas('point_types', [
            'code' => 'SP',
            'language' => 'Go',
            'label' => 'Go Skill Point',
        ]);
    }
}
