# Инструкция по добавлению логотипов подписок

Этот документ содержит список всех сервисов и инструкции по добавлению их официальных логотипов.

## Структура

Логотипы должны быть размещены в папке `/web/static/images/logos/` с названием файла в формате `service-name.png` или `service-name.svg`.

Рекомендуемый размер: 128x128px или 256x256px для PNG, векторный формат для SVG.

## Маркетплейсные сервисы

| Сервис | Официальный сайт | Имя файла |
|--------|------------------|-----------|
| MPSTATS | https://mpstats.io | mpstats.png |
| Маяк | https://mayak.chat | mayak.png |
| MarketGuru | https://marketguru.io | marketguru.png |
| EGGHEADS | https://eggheads.solutions | eggheads.png |
| Wildbox | https://wildbox.ru | wildbox.png |
| Stat4Market | https://stat4market.ru | stat4market.png |
| SellerStats | https://sellerstats.ru | sellerstats.png |
| woysa.club | https://woysa.club | woysa.png |
| XWAY | https://xway.pro | xway.png |
| Sellmonitor | https://sellmonitor.ru | sellmonitor.png |
| WBStat.PRO | https://wbstat.pro | wbstat.png |
| SellerFox | https://sellerfox.ru | sellerfox.png |
| Topseller | https://topseller.pro | topseller.png |
| Pi-Data | https://pi-data.ru | pidata.png |
| SelSup | https://selsup.ru | selsup.png |
| Модульселлер | https://moduleseller.ru | moduleseller.png |
| Salist | https://salist.ru | salist.png |
| Тинькофф Селлер | https://seller.tinkoff.ru | tinkoff-seller.png |
| ADAPTER | https://adapter.ru | adapter.png |
| MPSPACE | https://mpspace.ru | mpspace.png |
| Huckster | https://huckster.io | huckster.png |
| MP SURF | https://mpsurf.ru | mpsurf.png |
| inSales | https://insales.ru | insales.png |
| SalesFinder | https://salesfinder.ru | salesfinder.png |
| MP Manager | https://mpmanager.ru | mpmanager.png |
| Smart Seller | https://smartseller.su | smartseller.png |
| Indeepa | https://indeepa.ru | indeepa.png |
| Точка маркетплейсы | https://tochka.com | tochka-marketplace.png |
| JVO | https://jvo.ru | jvo.png |
| MarketProvider | https://marketprovider.ru | marketprovider.png |
| Ракета | https://raketa.world | raketa.png |
| SoykaSoft | https://soykasoft.ru | soykasoft.png |
| Mpfinassist | https://mpfinassist.ru | mpfinassist.png |
| WB FIN | https://wbfin.ru | wbfin.png |
| Таблички | https://tablicki.com | tablicki.png |
| Нейромаркет | https://neuromarket.ru | neuromarket.png |
| Sellego | https://sellego.ru | sellego.png |
| ОТВЕТО | https://otveto.ai | otveto.png |
| Отвечумба | https://otvechumba.ru | otvechumba.png |
| ТОРГСТАТ | https://torgstat.ru | torgstat.png |
| MPfact | https://mpfact.ru | mpfact.png |
| TrueStats | https://truestats.ru | truestats.png |
| MP Profit | https://mpprofit.ru | mpprofit.png |
| RASK | https://rask.ai | rask.png |
| TACTICs | https://tactics.ru | tactics.png |
| Спикс | https://spiks.ru | spiks.png |
| Джем | https://jam.team | jam.png |

## CRM системы

| Сервис | Официальный сайт | Имя файла |
|--------|------------------|-----------|
| AMOCRM | https://amocrm.ru | amocrm.png |
| Мегаплан | https://megaplan.ru | megaplan.png |
| Битрикс24 | https://bitrix24.ru | bitrix24.png |

## Как добавить логотипы

### Способ 1: Скачивание вручную

1. Перейдите на официальный сайт сервиса
2. Найдите логотип (обычно в header или footer, либо в разделе "Пресс-кит" / "Бренд")
3. Скачайте логотип в формате PNG или SVG
4. Переименуйте файл согласно таблице выше
5. Разместите файл в `/web/static/images/logos/`
6. Добавьте `logoUrl` в соответствующую запись в `web/static/js/app.js`

### Способ 2: Использование скрипта (для разработчиков)

Вы можете создать автоматический скрипт для загрузки логотипов:

```bash
#!/bin/bash
cd web/static/images/logos

# Пример для одного сервиса
curl -o mpstats.png "https://mpstats.io/logo.png"

# Повторите для каждого сервиса
```

## Обновление кода

После размещения логотипов обновите `web/static/js/app.js`:

```javascript
// Пример
{
    name: 'MPSTATS',
    logo: '📊',
    logoUrl: '/static/images/logos/mpstats.png',  // Добавить эту строку
    category: 'marketplace',
    price: 5990,
    currency: 'RUB',
    period: 'monthly',
    popular: true
},
```

## Альтернатива: Использование внешних URL

Вместо локального хранения можно использовать CDN или прямые ссылки:

```javascript
{
    name: 'MPSTATS',
    logo: '📊',
    logoUrl: 'https://mpstats.io/favicon.ico',  // Прямая ссылка
    // ...
}
```

**Примечание:** При использовании внешних URL убедитесь, что они доступны и не блокируются CORS.

## Рекомендации

1. **Формат**: Предпочтительно SVG для масштабируемости, PNG для простоты
2. **Размер**: 128x128px или 256x256px
3. **Фон**: Прозрачный (для PNG) или белый (для JPEG)
4. **Оптимизация**: Используйте https://tinypng.com/ для сжатия PNG
5. **Naming**: Используйте kebab-case (lowercase с дефисами)

## Fallback

Если логотип не загрузится, автоматически отобразится эмодзи из поля `logo`. Это обеспечивает graceful degradation.
