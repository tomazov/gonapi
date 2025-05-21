# 🧳 Бізнес-план: Агрегатор туристичних пропозицій Otpusk API 4.0

> **Проєкт базується на:**
> [github.com/koddr/tutorial-go-fiber-rest-api](https://github.com/koddr/tutorial-go-fiber-rest-api)
> **Ініціатор:** Євгеній Томазов
> **Дата:** 2025-04-07

---

## 📌 1. Призначення системи

"Агрегатор турів" — це масштабована система збору та обробки туристичних пропозицій з API різних туроператорів (ТО), яка забезпечує:

- централізований пошук по TO з різною логікою API (JSON, XML, SOAP);
- швидке кешування відповідей (Memcached);
- асинхронну обробку задач (RabbitMQ);
- модульну обробку для кожного оператора (через vendor-адаптери);
- масштабовану REST API для клієнтів.

---

## 🧱 2. Архітектура

### 🎯 Компоненти:

| Компонент           | Технологія             | Призначення                              |
|---------------------|------------------------|-------------------------------------------|
| REST API            | Go + Fiber             | Приймає запити `/getResults?...`         |
| Менеджер черг       | RabbitMQ               | Асинхронно виконує задачі                 |
| Тимчасовий кеш      | Memcached              | Зберігає статуси та результати (15 хв)   |
| БД операторів       | MySQL                  | Зберігає список ТО та параметри доступу  |
| Репозиторій цін     | ClickHouse             | Зберігання актуальних цін                 |
| Логи                | Файлова система        | `/var/log/vendor_<rec_id>.log`           |

---

## 🧠 3. Ключові механіки

### 🔑 `searchId` генерація

- Формується з параметрів запиту (`from`, `to`, `checkIn`, `food`, ...).
- Перетворюється в унікальний `md5` → `UUID` формат.
- Визначає унікальність пошуку → кеш + ідентифікація черги.

### 🧩 Memcached ключі:

| Ключ                                | Опис                                | TTL       |
|-------------------------------------|--------------------------------------|-----------|
| `search:<searchId>:workProgress`    | JSON зі статусом по всім TO          | 15 хв     |
| `search:<searchId>:<rec_id>:status` | Статус TO (`run`, `work`, `done`)    | 15 хв     |
| `search:<searchId>:<rec_id>:data`   | Результат обробки (JSON з цінами)    | 15 хв     |

---

### 🧵 Обробка Worker-ом:

1. Задача отримується з RabbitMQ.
2. Визначається `rec_id` → запускається відповідний handler.
3. Виконується API-запит до туроператора.
4. Результат кешується в `Memcached`.
5. Записується лог: `run` → `work` → `done`.

---

### 🪵 Логи:

- Лог-файл на кожного оператора:
  /var/log/vendor_<rec_id>.log
- Формат:

```
[2025-04-07 12:05:00] [9821] searchId=8dc99... operator=TPG2 status=run [2025-04-07 12:06:10] [9821] searchId=8dc99... operator=TPG2 status=done offers=124 duration=1.45s
```

---

## 🧩 4. Vendors: система адаптерів

### 🔌 Реєстрація адаптерів:

```
go
func init() {
  Register(2700, TPGHandler) // для TPG
}
```

### 📂 Приклади файлів:

```
/internal/vendors/samosoftApi.go
/internal/vendors/merlinxApi.go
/internal/vendors/coralTravelApi.go
/internal/vendors/proxymoPackagesApi.go
/internal/vendors/masterTourApi.go
```

### 🧠 Автоматична класифікація ТО по fApi:

```
if strings.Contains(fApi, "merlinx.eu") {
    // merlinxApi
} else if strings.Contains(fApi, "obs.md") {
    // tocoApi
} else if strings.Contains(fApi, "samo.") || fClassName IN ("JoinUp", "Pilon") {
    // samosoftApi
}
```
---

## 🧠 5. Формат відповіді API

```
{
  "lastResult": false,
  "results": {
    "2700": { "status": "done", "offers": [...] },
    "3306": { "status": "work", "offers": null }
  }
}
```

---

## ⚙️ 6. Масштабування

- Worker-обробка — в goroutines
- RabbitMQ – рівномірно розподіляє задачі
- Memcached – централізований TTL кеш
- Нові ТО додаються динамічно через Register(rec_id, handler)

---

## 📈 7. Економіка
- Перевага	Результат
- Асинхронність	Швидка відповідь без блокування
- Кешування	Зниження навантаження на зовнішні API
- Легке підключення нових ТО	Через vendors/*.go файли
- Висока продуктивність	1000+ обробок в секунду

---

## ✅ 8. Поточний статус

- ✅ Архітектура узгоджена
- ✅ TO список структуровано
- ✅ Реалізовано кеш + логіка TTL
- ⏳ В процесі: vendors/*.go + ClickHouse інтеграція

---

## 📌 9. Наступні кроки

- Реалізувати registry.go + Register() логіку
- Підключити всі активні ТО
- Інтегрувати ClickHouse запис цін
- Побудувати REST-ендпоінти в Fiber
- Розгорнути тестовий RabbitMQ кластер

---

## 🛠️ 10. Install GO

```
bash
cd ~
curl -O https://dl.google.com/go/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
sudo chown -R root:root /usr/local/go
sudo export PATH=$PATH:/usr/local/go/bin

cd /var/www/tomazov/tomazov.napi.otpusk
rm -rf ./main && /usr/local/go/bin/go build ./main.go

sudo cp ./napi.service /lib/systemd/system/napi.service
sudo chown -R root:root /lib/systemd/system/napi.service

sudo service napi stop
sudo service napi start
sudo service napi status

sudo cp ./napi.conf /etc/nginx/sites-available/napi.conf
sudo chown -R root:root /etc/nginx/sites-available/napi.conf

nginx -t
service nginx reload
```

---

# Завдання. Створення програми-агрегатора цін на Go, Використовуємо:

## 1. RabbitMQ — для запуска черг пошуку

### 1.1. приклад пошукового запита якій запускає пошук https://api.otpusk.com/getResults?from=1831&to=115&stars=4,5&checkIn=2024-10-01&checkTo=2024-10-08&nights=7&nightsTo=8&people=10709&food=uai,ai&transport=air&price=100&priceTo=5000&page=1&currencyLocal=eur&toOperators=3377,3374,3411,3441&availableFlight=yes,request,no&stopSale=yes,request,no&lang=ukr&group=1&rating=7-10&number=0&data=extlinks&access_token=29ae6-32ef8-d106b-7ea88-5f40d якщо його ще раз запитати він повинени вивести данні якщо він їх отримав від черги

### 1.2. опис отриманих даних

```
- from:1831 - місто відправлення
- to:115 - номер країни (до 500), міста (від 500 до 5000) чи готелю (від 5000)
- stars:4,5 - зірки готелів яки треба шукати
- checkIn:2024-10-01 - дата відпочинку шукаємо "від"
- checkTo:2024-10-08 - дата відпочинку шукаємо "до"
- nights:7 - кількість ночей "від"
- nightsTo:8 - кількість ночей "до"
- people:10709 - кількість туристів, де перша цифра це кількість дорослих (від 1 до 9), наступні десятки це вік дітей, у цьому прикладі це дві дитини 7 і 9 років
- food:uai,ai - шукаємо тип харчування
- transport:air - тип транспорту
- price:100 - ціна "від"
- priceTo:5000 - ціна "до"
- currency:EUR - показуємо ціни у валюті EUR
- page:1 - сторінка ціни
- currencyLocal:UAH - показуємо ціни у валюті EUR
- toOperators:3377,3374,3411,3441 - шукаємо лише в цих ТурОператорів
- availableFlight:yes,request,no - остановки продаж білетів в транспорті
- stopSale:yes,request,no - остановки продаж у готелі
- lang:UK - мова відображення текстів ISO 639-1
- group:1 - групування цін (можливо опишу далі)
- rating:7-10 - рейтинг готелів з бази даних
- number:0 - нумерація запита (максимум до 20 запитів)
- access_token:29ae6-32ef8-d106b-7ea88-5f40d - токен клієнта з бази даних
```

## 2. Memcached — для зберігання статусу обробки завдань у RabbitMQ і зберігання оброблених результатів 15 хвилин

### 2.1. на основі даних з пункту 1.2 генеруємо пошуковий ключ, приклад на PHP "$charId = md5($searcString); $searchId = substr($charId, 0, 8) . '-' . substr($charId, 8, 4) . '-' . substr($charId, 12, 4) . '-' . substr($charId, 16, 4) . '-' . substr($charId, 20, 12);". MySQL — для зберігання синхронізованих готелів та іншої структурованої інформації

## 3.ClickHouse — для зберігання цін на 5 діб, структура:

```
- operatorId — номер джерела ціни, номер ТурОператор
- countryId — номер країни відпочинку (від 0 до 500)
- cityId — номер міста відпочинку (від 500 до 5000)
- hotelId — номер готелю відпочинку (від 5000)
- fromCityId — номер міста відправлення
- adultChild — кількість туристів, наприклад 2 - це два туриста, 9040614 - це 9 туристів плюс дитина 4 років плюс дитина 6 років плюс дитина 14
- ages — діапазон дітей від до, якщо є діти, наприклад для 3 дітей такі діапазони [[0,4],[2,8],[6,16]] яки закодовані в бінарний код, ось приклад на PHP "function ages2bin($arrInputAges) { $arrAges = array(); for ($i = 0; $i <= 63; $i++) $arrAges[$i] = 0; for ($i = 0; $i <= 2; $i++) { if (!isset($arrInputAges[$i])) $arrInputAges[$i][0] = $arrInputAges[$i][1] = -1; if (is_array($arrInputAges[$i]) && $arrInputAges[$i][0] >= 0 && $arrInputAges[$i][1] >= 0 && $arrInputAges[$i][0] <= $arrInputAges[$i][1]) { or ($age = $arrInputAges[$i][0]; $age <= $arrInputAges[$i][1]; $age++) { $arrAges[19 * (2 - $i) + $age] = 1; } } return implode('', array_reverse($arrAges)); }"
- checkIn — дата початку відпочинку
- nights — тривалість подорожі
- duration — тривалість відпочинку в готелі
- transportFood — транспорт і харчування закодовано у двозначну цифру (десятки це транспорт (0:'no'|1:'air'|2:'bus'|3:'train'|4:'trainbus'|5:'ship'), одноцифрові це тип харчування (0:'ob'|1:'bb'|2:'hb'|3:'fb'|4:'ai'|5:'uai')), тобто 14 це транспорт='air' і харчування='ai'
- transportOption — транспорт
- tourName — назва туру
- stopSale — остановки продаж {"avia":("yes"|"request"|"no"),"aviaBack:("yes"|"request"|"no"),"hotel:("yes"|"request"|"no")}
- roomName — назва номеру
- price — ціна (для удобства зберігання в INT32 її помножаємо на 100) потів при відображеня поділімо на 100
- currency — валюта ціни, цифровий код (ISO 4217)
- tourOptions — допоміжні дані про ціну [1:"insurance",2:"transfer",3:"luggage"]
- bronUrl — посилання на бронювання
- offerId — номер ціни, закодованно віповідно до више отриманних данних: crc32(operatorId,'h',hotelId,'c',fromCity,'a',adultChild,'g',ages,'d',checkIn,'n',nights,'f',food,'r',crc32(roomName),'t',transport,'s',stopSale,'o',crc32(tourName),'o',tourOptions)
- active — актуальна ціна чи ні (1|0)
- updateTime — дата отримання ціни
```

## 4. для кожної ціни ми будемо генерувати offerID який будемо зберігати у clickhouse, очь аглоритм генерації на PHP

```
function genKey($offer, $hotelId, $childName = ''){
   $foods = array('ob' => 0, 'bb' => 1, 'hb' => 2, 'fb' => 3, 'ai' => 4, 'uai' => 5);
   $transports = array('no' => 1, 'air' => 2, 'bus' => 3, 'train' => 4, 'trainbus' => 5, 'ship' => 6);
      $offer->checkIn = trim($offer->checkIn);
      $offer->type = trim($offer->type);
      $offer->food = trim($offer->food);
      $offer->transport = trim($offer->transport);
      $val = array(
          (int)$offer->operatorId,
          'h',
          (int)$hotelId,
          'c',
          (int)$offer->fromCity,
          'd',
          str_replace('-', '', $offer->checkIn),
          (string)$offer->type,
          $offer->operatorId==3452? (int)$offer->duration: (int)$offer->length,
          (string)$offer->food,
          (int)$offer->roomId,
          (string)$offer->transport,
          (int)$offer->tourId,
      );
      if (!empty($offer->foodName) && $offer->food != $offer->foodName) {
          $food = preg_replace('![^A-Za-z\d-]!i', '', $offer->foodName);
          if ($food && $offer->food != $food) $val[] = intval(Service_Base62::crc8($food));
      }
      $childName = trim($childName);
      if ($offer->child > 0 && !empty($childName)) {
          $accName = trim(preg_replace('!^.*CHD!', '', $childName));
          $accName = str_replace('.00', '', $accName);
          $accName = str_replace('.99', '', $accName);
          $accName = preg_replace('![()]!', '', $accName);
          $val[] = 'a' . preg_replace('![^A-Za-z\d-]!i', '', mb_strtolower($accName));
      } elseif (!empty($childName)) {
          $val[] = 'a' . preg_replace('![^A-Za-z\d-]!i', '', mb_strtolower($childName));
      }
      $tourStatus = trim($offer->tourStatus);
      if (!empty($tourStatus)) {
          $val[] = preg_replace('![^A-Za-z\d]!i', '', mb_strtolower($tourStatus, 'UTF-8'));
      }
      $str = implode('', $val);
      // если перешли границу года, то год берем предыдущий
      $year = date('Y') - (date('Y', strtotime($offer->checkIn)) >= date('Y') ? 0 : 1);
      $day = date_diff(date_create($year . '-01-01'), date_create($offer->checkIn))->format('%a');
      $num = sprintf('%u', crc32($str));

      if (strlen($num) < 10) $num = str_pad($num, 10, 0);

      // добавляем к хешу уникальность в виде дня и длительности
      $num = substr_replace($num, sprintf('%03d', $day), 3, 0);
      $num = substr_replace($num, sprintf('%02d', $offer->length), -4, 0);
      $num = substr_replace($num, $foods[$offer->food], -2, 0);

      if ($num > 18446744073709551615) return null; // UInt64

      return $num;
   }
```

## 5. searchId - буде виравнюватись так і генеруватись

```
function _toString()   {
    $s = ($this->project == 'onsite' ? $this->project . '|' : '');
    $s .= 'f' . (int)$this->fromCity;
    $s .= (int)$this->id;
    $s .= 'd' . $this->checkIn;
    $s .= 'd' . $this->checkInTo;
    $s .= 'l' . (int)$this->length;
    $s .= 'l' . (int)$this->lengthTo;
    $s .= 'e' . (int)$this->people;
    $s .= 'p' . (int)$this->page;
    if (isset($this->option['oldPrice'])) {
        $s .= 'o';
    }
    if (in_array($this->sort, array('rating'))) {
        $s .= 's' . $this->sort;
    }
    if ($this->id < 5000 && !empty($this->group) && $this->group > 0 && $this->group < 256) {
        $s .= 'g' . $this->group;
    }
    if (!empty($this->vendor)) {
        $s .= 'v' . $this->vendor; // . (!empty($this->shard) ? '|' . $this->shard : '');
    } else {
        if (!empty($this->incVendors)) {
            $incVendors = explode(',', $this->incVendors);
            $s .= 'w' . count($incVendors);
        }
    }
    if ($this->directFlight) $s .= 'd1';
    if (!empty($this->transport) && is_array($this->transport) && count($this->transport) > 0 && count($this->transport) < 4) {
        // нормализируем данные
        $sort = $this->transport;
        sort($sort);
        $s .= 't' . implode(',', $sort);
    }
    if (!empty($this->food) && is_array($this->food) && count($this->food) > 0 && count($this->food) < 6) {
        // нормализируем данные
        $sort = $this->food;
        sort($sort);
        $s .= 'f' . implode(',', $sort);
    }
    // игнорируем другие данные #7243
    if (!empty($this->toHotels) && is_array($this->toHotels) && count($this->toHotels) > 0 && count($this->toHotels) <= 100) {
        // нормализируем данные
        $sort = array();
        foreach ($this->toHotels as $hid) {
            if ($hid >= 5000) {
                $sort[] = (int)$hid;
            }
        }
        if (count($sort)) {
            sort($sort);
            $s .= 'h' . implode(',', $sort);
        }
    } else {
        if (!empty($this->toCities) && is_array($this->toCities) && count($this->toCities) > 0 && count($this->toCities) <= 50) {
            // нормализируем данные
            $sort = array();
            foreach ($this->toCities as $cid) {
                if ($cid >= 500 && $cid < 5000) {
                    $sort[] = (int)$cid;
                }
            }
            // проверям выбраны ли все курорты по стране с учетом цены
            if (count($sort) > 2) {
                $cities = DB::select('rec_id')
                    ->from('priceIndex')
                    ->where('rec_id', '>', 500)
                    ->where('rec_id', '<', 5000)
                    ->where('rec_id', 'IN', $sort)
                    ->where('price', '>', 0)
                    ->where('countryId', '=', $this->store['id'])
                    ->group_by('rec_id')
                    ->cached()->execute('otpusk')->as_array(null, 'rec_id');
            } else {
                $cities = $sort;
            }
            if (count($sort) && count($cities)) {
                sort($sort);
                $s .= 's' . implode(',', $sort);
            }
            unset($cities);
        }
        if (!empty($this->toProvinces) && is_array($this->toProvinces) && count($this->toProvinces) > 0 && count($this->toProvinces) <= 50) {
            // нормализируем данные
            $sort = $this->toProvinces;
            if (count($sort)) {
                sort($sort);
                $s .= 'p' . implode(',', $sort);
            }
        }
        if (!empty($this->toDistricts) && is_array($this->toDistricts) && count($this->toDistricts) > 0 && count($this->toDistricts) <= 50) {
            // нормализируем данные
            $sort = $this->toDistricts;
            if (count($sort)) {
                sort($sort);
                $s .= 'd' . implode(',', $sort);
            }
        }
        if (!empty($this->stars) && is_array($this->stars) && count($this->stars) > 0 && count($this->stars) < 6) {
            // нормализируем данные
            $sort = $this->stars;
            sort($sort);
            $s .= 'r' . implode(',', $sort);
        }
        if (!empty($this->services) && is_array($this->services) && count($this->services) > 0) {
            // нормализируем данные
            $sort = $this->services;
            sort($sort);
            $s .= 'i' . implode(',', $sort);
        }
        if (!empty($this->ignoreServices) && is_array($this->ignoreServices) && count($this->ignoreServices) > 0) {
            // нормализируем данные
            $sort = $this->ignoreServices;
            sort($sort);
            $s .= 'q' . implode(',', $sort);
        }
        if (!empty($this->rate) && is_array($this->rate) && count($this->rate) > 0) {
            if (!(count($this->rate) == 2 && $this->rate[0] == 0 && $this->rate[1] == 10)) {
                // нормализируем данные
                $sort = $this->rate;
                sort($sort);
                $s .= 'r' . implode(',', $sort);
            }
        }
        if (!empty($this->rating) && is_array($this->rating) && count($this->rating) > 0) {
            if (!(count($this->rating) == 2 && $this->rating[0] == 0 && $this->rating[1] == 10)) {
                // нормализируем данные
                $sort = $this->rating;
                sort($sort);
                $s .= 'p' . implode(',', $sort);
            }
        }
        if (!empty($this->cat) && is_array($this->cat) && count($this->cat) > 0) {
            // нормализируем данные
            $sort = $this->cat;
            sort($sort);
            $s .= 'c' . implode(',', $sort);
        }
        if (!empty($this->form) && is_array($this->form) && count($this->form) > 0) {
            // нормализируем данные
            $sort = $this->form;
            sort($sort);
            $s .= 'f' . implode(',', $sort);
        }
        if (!empty($this->reviews)) {
            $s .= 'w' . (int)$this->reviews;
        }
        if (!empty($this->chainId) && is_array($this->chainId) && count($this->chainId) > 0) {
            $s .= 'n' . implode(',', $this->chainId);
        }
    }
    if (!empty($this->toNotHotels) && is_array($this->toNotHotels) && count($this->toNotHotels) > 0) {
        // нормализируем данные
        $sort = array();
        foreach ($this->toNotHotels as $hid) {
            if ($hid >= 5000) {
                $sort[] = (int)$hid;
            }
        }
        if (count($sort)) {
            sort($sort);
            $s .= 'n' . implode(',', $sort);
        }
    }
    if (!empty($this->currency) && $this->currency != 'uah' && $this->currency != $this->currencyLocal) {
        $s .= 'u' . strtolower((string)$this->currency);
    }
    if (!empty($this->currencyLocal) && ($this->currencyLocal != 'uah' || $this->currencyLocal != $this->currency)) {
        $s .= 'ul' . strtolower((string)$this->currencyLocal);
    }
    if (!empty($this->minPrice)) {
        if (empty($this->currency) || $this->currency == 'uah') {
            $s .= 'n' . round($this->minPrice / 100) * 100;
        } else {
            $s .= 'n' . round($this->minPrice / 10) * 10;
        }
    }
    if (!empty($this->maxPrice)) {
        if (empty($this->currency) || $this->currency == 'uah') {
            $s .= 'x' . round($this->maxPrice / 100) * 100;
        } else {
            $s .= 'x' . round($this->maxPrice / 10) * 10;
        }
    }
    if (!empty($this->availableFlight)) {
        $sort = $this->availableFlight;
        sort($sort);
        $s .= 'a' . implode(',', $sort);
    }
    if (!empty($this->stopSale)) {
        $sort = $this->stopSale;
        sort($sort);
        $s .= 's' . implode(',', $sort);
    }
    if (!empty($this->toFilters) && is_array($this->toFilters) && count($this->toFilters) > 0) {
        // нормализируем данные
        $fkey = array();
        foreach ($this->toFilters as $dir => $geos) {
            foreach ($geos as $geo => $data) {
                $fkey[] = dechex($dir) . '.' . dechex($geo);
            }
        }
        $s .= 'f' . implode(',', $fkey);
        unset($fkey);
    }
    // Для генерации основного ключа задания и для отпуска
    if (empty($this->vendor) || (!empty($this->vendor) && $this->vendor === 'otpk')) {
        if (!empty($this->toOperators) && is_array($this->toOperators) && count($this->toOperators) > 0) {
            // нормализируем данные
            $sort = $this->toOperators;
            sort($sort);
            $s .= 'o' . implode(',', $sort);
        }
    }
    if ($this->autoDirections != 'yes' && !empty($this->toDirections) && is_array($this->toDirections)
        && count(
            $this->toDirections
        ) > 0
    ) {
        // только для ДБ отпуска
        if (!empty($this->vendor) && $this->vendor === 'otpk') {
            // нормализируем данные
            $sort = $this->toDirections;
            sort($sort);
            $s .= 'd' . implode(',', $sort);
        }
    } elseif (!empty($this->debug) && $this->debug == 99) {
        $s .= 'd' . rand(1000000, 9999999);
    }
    if ($this->roomId > 1) {
        $s .= 'r'.$this->roomId;
    }
    if (!empty($this->returnNow)) {
        //$s .= 'returnNow';
    }
    if (!empty($this->noPromo)) {
        $s .= 'noPromo';
    }
    if (!empty($this->offers) && $this->vendor !== 'otpk') {
        $s .= 'offers';
    } elseif (!empty($this->offerId) && $this->vendor !== 'otpk') {
        $s .= 'f' . (int)$this->offerId;
    }
    if (!empty($this->callNext)) {
        $s .= 'callNext';
    }
    if (!empty($this->version) && $this->version >= 2.5) {
        $s .= 'v'.str_replace('.','',$this->version);
    }
    if (!empty($this->multiSearch) && strlen($this->multiSearch) > 1) {
        $s .= 'm'.$this->multiSearch;
    }
    if ($this->limit > 10) {
        $s .= 'l'.$this->limit;
    }
    return $s;
}
```

## 6. дивись Samo кліентів може буде багато більше 20 штук, треба буде брати посилання з бази mysql `SELECT rec_id,fClassName,fApi,fToken,fProxy,fTimeOut,fUserLogin,fUserPass FROM `tOperators` WHERE fActive='yes' AND fMain='no' and fClassName!=''` rec_id - це номер оператора, він треба буде для ініціалізацію класу і записів у лог для цього ТО, fApi - це куди робити запит, fClassName - це назва классу (якщо потрібно), fProxy - ролити запити через проксі, fTimeOut - скільки часу чекати відповіді, fToken - доступ до апіf, UserLogin,fUserPass  - доступ до апі (якщо буде необхідно)
