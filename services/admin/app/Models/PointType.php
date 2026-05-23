<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class PointType extends Model
{
    protected $table = 'point_types';

    protected $primaryKey = 'code';

    public $incrementing = false;

    public $timestamps = false;

    protected $fillable = ['code', 'language', 'label'];

    protected $casts = [];

    protected function setKeysForSaveQuery($query)
    {
        return $query
            ->where('code', '=', $this->original['code'] ?? $this->getAttribute('code'))
            ->where('language', '=', $this->original['language'] ?? $this->getAttribute('language'));
    }
}
