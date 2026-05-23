<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class PointType extends Model
{
    protected $table = 'point_types';
    protected $primaryKey = null;
    public $incrementing = false;
    public $timestamps = false;

    // code・language は作成後 immutable なので fillable に含めない
    protected $fillable = ['label'];

    protected $casts = [];
}
