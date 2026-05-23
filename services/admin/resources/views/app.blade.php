<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{ config('app.name', 'Lang War Admin') }}</title>
    @inertiaHead
    @vite(['resources/js/app.ts'])
</head>
<body class="bg-gray-100">
    @inertia
</body>
</html>
