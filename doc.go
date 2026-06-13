/*
Пакет res — хранилище ресурсов приложения с поддержкой Transform и интеграции с sdi.

Сборка ресурсов (config, builder) выполняется снаружи; res принимает готовые объекты
через AddBuiltin (system) или Add / AddAll (user, builder.Registrar).

Devconv:
  - AddBuiltin — system defaults из init (import _ "…/logger")
  - Add — user config из builder.Build
  - enforcement — golangci profiles в github.com/omcrgnt/lint

Основные возможности:
  - AddBuiltin / Add / AddAll / Default.Add: регистрация ресурсов
  - Origin (System | User), WalkEntries — metadata для Dedup
  - Dedup — interface dedupe executor для sdi.Resolve (policy в sdi)
  - Transform: подготовка ресурсов перед wiring
  - Default / Walk: read-only pool для sdi.Resolve
  - Get / Find: типизированный доступ

Типичный pipeline:

	ecfg.Parse → builder.Build(cfg, res.Default)
	res.Transform(obs.Instrument)
	sdi.Resolve(res.Default)

Ограничения:
  - один concrete type — один ресурс в byType
  - Add заменяет system-ресурс того же concrete type
  - несколько implementor одного interface допустимы до sdi.Resolve
  - Transform выполнять до sdi.Resolve
*/
package res
