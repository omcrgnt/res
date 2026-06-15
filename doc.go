/*
Package res — регистрация и хранение ресурсов приложения.

Центральный тип — [Registry] ([Default], [New]): add, query, transform, remove.
Ресурс в pool представлен [Entry] (concrete type, value, опциональные [Tag]).
Один type может иметь несколько [Entry]; [Registry.GetByType] и
[Registry.GetByInterface] возвращают все совпадения (0..N),
[Registry.GetOneByType] и [Registry.GetOneByInterface] — первое в порядке
регистрации или [ErrNotFound].

[Tag] — метаданные на [Entry], задаются через [Registry.AddWithTags].
[Registry] теги только хранит и отдаёт в [Entry.Has]; сам по ним не действует.
[TagReplaceable] — «запасной вариант»: caller при выборе одного ресурса из
нескольких может предпочесть запись без этого тега и убрать остальные через
[Registry.Remove].

API ([Default] и package-level funcs):
  - create: [Add], [AddWithTags]
  - read: [WalkEntries], [GetByType], [GetByInterface], [GetOneByType], [GetOneByInterface]
  - update: [Transform]
  - delete: [Remove]

Devconv:
  - [AddWithTags] с [TagReplaceable] — library defaults (import _ "…/pkg")
  - [Add] — явная регистрация caller'ом
  - enforcement — golangci profiles в github.com/omcrgnt/lint
*/
package res
