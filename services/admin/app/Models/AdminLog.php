<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class AdminLog extends Model
{
    public $timestamps = false;

    protected $fillable = [
        'action',
        'target_type',
        'target_id',
        'payload',
        'created_at'
    ];

    protected $casts = [
        'payload'    => 'array',
        'created_at' => 'datetime',
    ];

    public static function record(
        string $action,
        string $targetType,
        string $targetId,
        array $payload
    ): void
    {
        static::create([
            'action'      => $action,
            'target_type' => $targetType,
            'target_id'   => $targetId,
            'payload'     => $payload,
            'created_at'  => now(),
        ]);
    }
}
