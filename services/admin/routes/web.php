<?php

use App\Http\Controllers\GuildController;
use App\Http\Controllers\PointTypeController;
use App\Http\Middleware\BasicAuthMiddleware;
use Illuminate\Support\Facades\Route;
use Inertia\Inertia;

Route::middleware(BasicAuthMiddleware::class)->group(function () {
    Route::get('/', fn () => redirect()->route('guilds.index'));

    Route::get('/guilds', [GuildController::class, 'index'])->name('guilds.index');
    Route::get('/guilds/create', [GuildController::class, 'create'])->name('guilds.create');
    Route::post('/guilds', [GuildController::class, 'store'])->name('guilds.store');
    Route::get('/guilds/{guild}/edit', [GuildController::class, 'edit'])->name('guilds.edit');
    Route::put('/guilds/{guild}', [GuildController::class, 'update'])->name('guilds.update');

    Route::get('/point-types', [PointTypeController::class, 'index'])->name('point-types.index');
    Route::get('/point-types/create', [PointTypeController::class, 'create'])->name('point-types.create');
    Route::post('/point-types', [PointTypeController::class, 'store'])->name('point-types.store');
    Route::get('/point-types/edit', [PointTypeController::class, 'edit'])->name('point-types.edit');
    Route::put('/point-types', [PointTypeController::class, 'update'])->name('point-types.update');
});
