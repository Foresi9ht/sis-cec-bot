package complaint

// Keyword lists used by KeywordProcessor, in Russian, English and Kazakh.
//
// Matching rules (see matcher.go):
//   - the text is lower-cased and split into words;
//   - a keyword matches when a word starts with it, so "нарушен" finds
//     "нарушение" and "нарушения" (these are word stems);
//   - a keyword ending with "$" must match a whole word, which is used for
//     short words that would otherwise match too much ("суп$" must not match
//     "супер", "кран$" must not match "кранов");
//   - a keyword may contain several words ("качество еды").
//
// Duplicates are removed when the lists are compiled.

// validKeywords point to the CEC's jurisdiction: elections, rights violations,
// misconduct of staff and student organizations, administration.
var validKeywords = []string{
	// elections
	"выборы", "кандидат", "голосован", "агитац", "фальсификац",
	"election", "candidate", "voting", "agitation", "falsif",
	"сайлау", "үміткер", "дауыс", "үгіт",

	// rights and violations
	"нарушен", "давлен", "принужден", "угроз", "шантаж",
	"взятк", "коррупц", "дискриминац", "буллинг", "harassment",
	"violation", "pressure", "forced", "threat", "blackmail",
	"bribe", "corruption", "discriminat", "bullying",
	"бұзушылық", "қысым", "мәжбүр", "қауіп", "бопсалау",
	"пара$", "сыбайлас", "кемсіту",

	// teachers
	"преподаватель", "сотрудник", "оценивани", "занижен", "предвзятост",
	"субъективност", "некомпетент", "унижен", "агрессив", "навязыван",
	"дискриминац", "пересдач",
	"teacher", "staff", "grading", "bias", "incompetent",
	"humiliat", "aggressive", "retake",
	"оқытушы", "қызметкер", "бағалау", "төмендету", "біржақтылық",
	"кемсіту", "қорлау",

	// organizations
	"организац", "студ", "права", "устав", "регламент", "комитет",
	"руководств", "управлени", "игнор", "собрани",
	"organization", "student", "rights", "charter", "committee",
	"leadership", "management", "ignored", "meeting",
	"ұйым", "студент", "құқық", "жарғы", "комитет",
	"басшылық", "басқару", "елемеу",

	// administration
	"несправедлив", "бюрократ", "прозрачност", "дисциплинарн",
	"unfair", "bureauc", "transparent", "disciplin",
	"әділетсіз", "бюрократ", "ашықтық", "тәртіптік",
}

// invalidKeywords point to operational matters outside the CEC's jurisdiction:
// food, cleaning, plumbing, furniture, living conditions, technical services.
var invalidKeywords = []string{
	// food
	"еда$", "меню", "суп$", "обед", "завтрак", "ужин", "вкус", "готовят",
	"порци", "качество еды", "испорч", "волос в еде", "запах еды",
	"нездоров", "однообразн", "цена еды", "касса", "подсчёт", "подсчет",
	"food", "menu", "soup$", "lunch", "breakfast", "dinner", "taste",
	"portion", "food quality", "spoiled", "hair in food",
	"тамақ", "мәзір", "сорпа", "түскі ас", "таңғы ас", "кешкі ас",
	"бөлік", "тамақ сапасы",

	// cleaning
	"уборка", "грязь", "мусор", "антисанитар",
	"cleaning", "dirt", "trash", "unsanitary",
	"тазалық", "қоқыс",

	// plumbing
	"сантехник", "кран$", "душ$", "фильтр воды", "фонтанчик",
	"plumber", "faucet", "shower", "water filter",
	"су сүзгісі",

	// furniture
	"кровать", "матрас", "сломан дверь",
	"bed$", "mattress", "broken door",
	"төсек", "сынған есік",

	// living conditions
	"холод", "сквозняк", "шум$", "сосед", "разогрев еды",
	"cold$", "draft", "noise", "neighbor",
	"суық", "желаяқ", "шу$", "көрші",

	// technical
	"интернет", "wifi", "wi-fi", "свет$", "лампочка", "печать", "принтер",
	"internet", "light$", "bulb$", "printer",
	"жарық", "шам$",

	// activities (not violations)
	"скучн", "мало активност", "слабый промоушн",
	"boring", "few activities",
	"жалықтырарлық", "аз іс-шара",
}

// urgentKeywords raise the priority to HIGH: threats, violence, corruption,
// discrimination, falsification, food safety, unlawful eviction, smoking.
var urgentKeywords = []string{
	// physical threat
	"срочно", "немедленно", "угроза", "насилие", "принужден", "давлен",
	"шантаж", "взятка", "фальсификац", "критично", "опасно",
	"urgent", "immediately", "threat", "violence", "forced", "pressure",
	"blackmail", "bribe", "falsif", "critical", "dangerous",
	"шұғыл", "дереу", "қауіп", "зорлық", "мәжбүр", "қысым",
	"бопсалау", "пара$", "маңызды", "қауіпті",

	// humiliation
	"унижен", "оскорблен", "агрессив", "кричал", "угрожал",
	"humiliat", "insulted", "aggressive", "shouted", "threatened",
	"қорлау", "кемсіту", "айқайлады", "қорқытты",

	// corruption
	"коррупц", "взятк", "подкуп", "вымогательств",
	"corruption", "bribery", "extortion",
	"сыбайлас", "сатып алу", "қорқытып алу",

	// discrimination
	"дискриминац", "гендерн", "расизм", "ксенофоб",
	"discriminat", "gender", "racism", "xenophob",
	"гендер", "нәсілшілдік",

	// falsification
	"подделк", "сфабрик",
	"forgery", "fabricat",
	"бұрмалау", "жасанды",

	// food safety
	"испорч", "отравлен", "антисанитар", "волос в еде",
	"spoiled", "poisoned", "unsanitary", "hair in food",
	"бұзылған тамақ", "улану",

	// forced removal
	"выгоняют", "выселен", "незаконно выселил",
	"expelled", "evicted", "unlawfully evicted",
	"шығарып жіберу", "заңсыз шығару",

	// smoking
	"курен", "курит в здании", "курит внутри",
	"smoking inside", "smoking in building",
	"ғимарат ішінде темекі",
}

// mediumKeywords raise the priority to MEDIUM when no urgent keyword is found.
var mediumKeywords = []string{
	// violations
	"нарушен", "жалоб", "проблем", "некорректн", "несправедлив",
	"violation", "complaint", "problem", "incorrect", "unfair",
	"бұзушылық", "шағым", "мәселе", "дұрыс емес", "әділетсіз",

	// grading
	"занижен", "предвзятост", "субъективност", "критери оценивани",
	"пересдач", "изменен правил", "баллы", "ошибк в тест",
	"grading", "bias", "subjective", "grading criteria",
	"retake", "rules changed", "scores", "test error",
	"бағалау", "біржақтылық", "субъективтілік", "қайта тапсыру",
	"ереже өзгерту", "балл", "тест қатесі",

	// teaching
	"некомпетент", "плохое объяснен", "не готовится", "делегирован",
	"ассистент без проверки", "учебного плана", "структуры нет",
	"слабый уровень",
	"incompetent", "poor explanation", "unprepared", "delegated",
	"unverified assistant", "no syllabus", "no structure",
	"weak teaching",
	"біліксіз", "нашар түсіндіру", "дайындалмаған", "делегациялау",
	"оқу жоспары жоқ", "нашар оқыту",

	// organization
	"неактивн", "плохой менеджмент", "коммуникац", "игнор мнени",
	"хаос в чат", "несоответстви", "отсутств собрани", "вовлеченност",
	"inactive", "poor management", "communication", "ignored opinions",
	"chat chaos", "no meetings", "engagement",
	"белсенді емес", "нашар басқару", "байланыс", "пікірді елемеу",
	"жиналыс жоқ",

	// administration
	"несправедлив решен", "бюрократ", "прозрачност", "дисциплинарн",
	"add/drop", "сроки",
	"unfair decision", "bureauc", "transparent", "disciplin",
	"deadline",
	"әділетсіз шешім", "ашықтық", "тәртіптік",

	// behavior
	"оскорблен в чат", "нарушение правил",
	"insults in chat", "rule violation",
	"чатта қорлау", "ережені бұзу",

	// infrastructure
	"не работает wifi", "нет сети", "нет мест", "коворкинг",
	"техподдержк", "игнор заявк",
	"wifi down", "no network", "no seating", "coworking",
	"tech support", "ignored request",
	"wifi жоқ", "желі жоқ", "орын жоқ",

	// inclusivity
	"иностранн студент", "интеграц", "инклюзив",
	"international student", "integrat", "inclusiv",
	"шетелдік студент", "интеграция", "инклюзивтілік",
}
