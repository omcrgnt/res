/*
Пакет res — хранилище ресурсов приложения с поддержкой Transform и интеграции с sdi.

Сборка ресурсов (config, builder) выполняется снаружи; res принимает готовые объекты
через Add, AddAll или res.Default.Add (builder.Registrar).

Основные возможности:
  - Add / AddAll / Default.Add: регистрация готовых ресурсов
  - Transform: подготовка ресурсов перед wiring (obs, метрики и т.д.)
  - Default / Walk: read-only pool для sdi.Resolve
  - Get / Find: типизированный доступ после wiring

Типичный pipeline (оркестрация в main):

	ecfg.Parse → builder.Build(cfg, res.Default)
	res.Transform(obs.Instrument)
	sdi.Resolve(res.Default)
	res.Get / res.Find

Ограничения:
  - один concrete type — один ресурс в byType
  - Transform выполнять до sdi.Resolve
  - после Transform с обёрткой Get по старому concrete type может не сработать — используйте Find по интерфейсу
*/
package res
