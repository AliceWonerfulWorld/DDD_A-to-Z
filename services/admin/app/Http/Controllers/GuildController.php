<?php

namespace App\Http\Controllers;

use App\Models\AdminLog;
use App\Models\Guild;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Str;
use Inertia\Inertia;
use Inertia\Response;

class GuildController
{
    public function index(): Response
    {
        $guilds = Guild::orderBy('sort_order')->orderBy('name')->get();

        return Inertia::render('Guilds/Index', [
            'guilds' => $guilds,
        ]);
    }

    public function create(): Response
    {
        return Inertia::render('Guilds/Create');
    }

    public function store(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'slug'        => ['required', 'string', 'max:255', 'unique:guilds,slug', 'regex:/^[a-z0-9-]+$/'],
            'name'        => ['required', 'string', 'max:255'],
            'description' => ['required', 'string'],
            'icon'        => ['required', 'string', 'max:255'],
            'color'       => ['required', 'string', 'regex:/^#[0-9A-Fa-f]{6}$/'],
            'sort_order'  => ['required', 'integer', 'min:0'],
        ]);

        $guild = Guild::create(array_merge($validated, [
            'id'         => Str::ulid(),
            'created_at' => now(),
            'updated_at' => now(),
        ]));

        AdminLog::record('created', 'guild', $validated['slug'], $validated);

        return redirect()->route('guilds.index')->with('success', 'ギルドを作成しました');
    }

    public function edit(Guild $guild): Response
    {
        return Inertia::render('Guilds/Edit', [
            'guild' => $guild,
        ]);
    }

    public function update(Request $request, Guild $guild): RedirectResponse
    {
        $validated = $request->validate([
            'name'        => ['required', 'string', 'max:255'],
            'description' => ['required', 'string'],
            'icon'        => ['required', 'string', 'max:255'],
            'color'       => ['required', 'string', 'regex:/^#[0-9A-Fa-f]{6}$/'],
            'sort_order'  => ['required', 'integer', 'min:0'],
        ]);

        $guild->update($validated);

        AdminLog::record('updated', 'guild', $guild->slug, $validated);

        return redirect()->route('guilds.index')->with('success', 'ギルドを更新しました');
    }

    // 削除は意図的に実装しない（FK参照があるため）
}
