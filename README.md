# golang-conf-2026-samples

Код к докладу «Game Hacking на Go: что может язык там, где его не ждут» с [GolangConf 2026](https://golangconf.ru/2026/abstracts/17430).

Всё крутится вокруг GTA: San Andreas.

## cmd/hesoyam_decrease

Патчит память запущенной игры через `WriteProcessMemory`. После запуска HESOYAM начинает отнимать миллион вместо того, чтобы его давать.

## cmd/injector

Инжектор DLL: `VirtualAllocEx` + `WriteProcessMemory` + `CreateRemoteThread` на `LoadLibraryA`, потом второй `CreateRemoteThread` на экспорт из DLL.

```
injector -window "GTA: San Andreas" -dll cheat.dll -dll-proc GoMain
```

## cmd/new_cheatcode

DLL на Go (`-buildmode=c-shared`), которую инжектор кладёт в процесс. Регистрирует новый читкод `GNALOG` (GOLANG наоборот) и вызывает функции игры через `syscall.SyscallN` по абсолютным адресам.

Для сборки нужен кросс-компилятор mingw-w64 i686:

```
CC=i686-w64-mingw32-gcc GOOS=windows GOARCH=386 CGO_ENABLED=1 \
    go build -buildmode=c-shared -o cheat.dll ./cmd/new_cheatcode
```

## internal/winutils

Обёртки над WinAPI, которых не хватает в `golang.org/x/sys/windows`.
