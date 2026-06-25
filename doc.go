/*
Package res — регистрация и хранение ресурсов приложения.

Центральный тип — [Registry] ([Global], [New]): add, query, transform, remove.
Ресурс в pool представлен [Entry] (concrete type, value, опциональные [Tag]).
Один type может иметь несколько [Entry]; [Registry.GetByType] и
[Registry.GetByInterface] возвращают все совпадения (0..N),
[Registry.GetOneByType] и [Registry.GetOneByInterface] — первое в порядке
регистрации или [ErrNotFound].

Config-entry и resource-entry

До [github.com/omcrgnt/builder].Build в pool лежат config-entry — значения,
реализующие Build() (any, error) ([github.com/omcrgnt/builder].Builder).
Их регистрируют library use init через [MustAddToGlobalWithTags] (Replaceable/Fixed defaults)
и [github.com/omcrgnt/ecfg].Register из AppConfig ([Add], explicit, без tags).

Build обходит registry, вызывает Build(), регистрирует resource-entry с
наследованием [Entry.Tags] и удаляет config-entry через [Registry.Remove].
После Build в pool остаются resource-entry (и legacy non-Builder entries, если есть).

Pipeline:

	ecfg.Parse → ecfg.Register(cfg, reg) → builder.Build(reg) → reg.Transform → sdi.Resolve

[Tag] — метаданные на [Entry], задаются через [Registry.AddWithTags].
[Registry] теги только хранит и отдаёт в [Entry.Tags]; сам по ним не действует.
[TagReplaceable] — «запасной вариант»: caller при выборе одного ресурса из
нескольких может предпочесть запись без этого тега и убрать остальные через
[Registry.Remove] (интерпретирует [github.com/omcrgnt/sdi] при dedup).

[TagFixed] — «не подлежит подмене»: при 2+ entries одного dep type
[sdi.Resolve] завершается ошибкой; запись не удаляется dedup policy.

API — методы [Registry] на [Global] или на явном reg из [New]:

	reg.Add(v) — только [NewResourceer] или [BuildConfiger]
	reg.WalkEntries(fn)
	reg.Transform(...)

Library use init регистрирует defaults через [MustAddToGlobalWithTags] в [Global].

Subpackage [github.com/omcrgnt/res/restest] — optional test helpers over the same [Registry] API.
*/
package res
