<?php

namespace App\Http\Controllers;

use App\Models\AdminLog;
use App\Models\PointType;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\Rule;
use Inertia\Inertia;
use Inertia\Response;

class PointTypeController
{
    public function index(): Response
    {
        $pointTypes = PointType::orderBy('code')->orderBy('language')->get();

        return Inertia::render('PointTypes/Index', [
            'pointTypes' => $pointTypes,
        ]);
    }

    public function create(): Response
    {
        return Inertia::render('PointTypes/Create');
    }

    public function store(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'code'     => ['required', 'string', 'max:255'],
            'language' => ['required', 'string', 'max:255'],
            'label'    => ['required', 'string', 'max:255'],
        ]);

        // 複合PKの重複チェック
        $exists = PointType::where('code', $validated['code'])
            ->where('language', $validated['language'])
            ->exists();

        if ($exists) {
            return back()->withErrors(['code' => 'この code / language の組み合わせはすでに存在します']);
        }

        PointType::create($validated);

        AdminLog::record('created', 'point_type', "{$validated['code']}:{$validated['language']}", $validated);

        return redirect()->route('point-types.index')->with('success', 'ポイントタイプを作成しました');
    }

    public function edit(Request $request): Response
    {
        $pointType = PointType::where('code', $request->query('code'))
            ->where('language', $request->query('language'))
            ->firstOrFail();

        return Inertia::render('PointTypes/Edit', [
            'pointType' => $pointType,
        ]);
    }

    public function update(Request $request): RedirectResponse
    {
        $code     = $request->input('code');
        $language = $request->input('language');

        $pointType = PointType::where('code', $code)
            ->where('language', $language)
            ->firstOrFail();

        $validated = $request->validate([
            'label' => ['required', 'string', 'max:255'],
        ]);

        // code・language は immutable なので label のみ更新
        $pointType->update(['label' => $validated['label']]);

        AdminLog::record('updated', 'point_type', "{$code}:{$language}", $validated);

        return redirect()->route('point-types.index')->with('success', 'ポイントタイプを更新しました');
    }

    // 削除は意図的に実装しない（FK参照があるため）
}
