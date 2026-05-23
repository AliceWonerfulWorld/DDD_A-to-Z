<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Guild extends Model
{
    protected $table = 'guilds';
    protected $keyType = 'string';
    public $incrementing = false;

    // id・slug は作成後 immutable なので fillable に含めない
    protected $fillable = [
        'name',
        'description',
        'icon',
        'color',
        'sort_order',
    ];

    protected $casts = [
        'sort_order'  => 'integer',
        'current_exp' => 'integer',
        'guild_level' => 'integer',
        'created_at'  => 'datetime',
        'updated_at'  => 'datetime',
    ];
}
