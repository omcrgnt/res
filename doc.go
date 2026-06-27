/*
Package res — регистрация и хранение ресурсов приложения.

Центральный тип — [Registry] ([Global], [New]): add, query, transform, remove.
Ресурс в pool представлен [Entry] — handle на slot (Type, Value, system tags, custom tags).
Один type может иметь несколько [Entry]; [Registry.GetByType] и
[Registry.GetByInterface] возвращают все совпадения (0..N),
[Registry.GetOneByType] и [Registry.GetOneByInterface] — первое в порядке
регистрации или [ErrNotFound].

[Tag] — system metadata ([TagRegular], [TagReplaceable], [TagFixed]); задаётся через
[Registry.AddWithTags]. Интерпретация policy — в [github.com/omcrgnt/res/unique].

Custom tags — map[string]any на entry; registry хранит, не интерпретирует.
Чтение: [Entry.GetCustomTag]. Запись — через [github.com/omcrgnt/res/unique].AddWithCustomTag
или [AddWithTagsAndCustomTags] (для unique).

[Entry.ChangeValue] заменяет value in-place, сохраняя tags и custom tags.
Unique policy (one concrete type per registry) — в [github.com/omcrgnt/res/unique].

Type-unique registry (composition root): [github.com/omcrgnt/res/unique].

[Global] и [AddToGlobalWithTags] — legacy; prefer [unique.Global] и [unique.MustAddReplaceable].

Subpackage [github.com/omcrgnt/res/unique] — type-unique registry for app composition root.
Subpackage [github.com/omcrgnt/res/restest] — optional test helpers over the same [Registry] API.
*/
package res
