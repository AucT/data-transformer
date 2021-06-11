# data transformer

data transformer - командна утиліта написана на go, що дозволяє додати колонку до csv з іншого файлу.

## Базове використання
1. Витягти свої дані за допомогою [data-exporter](https://github.com/AucT/data-exporter)  
2. У файлі config.json змінити значення inputFileName, замість inputFileName.csv вписати ваш файл з оцінками, що ви отримали у першому пункті
3. Запустити data-transformer.exe (для windows) чи data-transformer (для linux)

В результаті ви отримаєте три файли
 - mappedOutputFileName.csv - таблиця тільки тих ваших оцінок, де скрипт знайшов ідентифікатор imdb
 - notMappedOutputFileName.csv - таблиця тільки тих ваших оцінок, де скрипт не знайшов ідентифікатор imdb
 - combinedOutputFileName.csv - таблиця всіх ваших оцінок. Не знайдені imdb будуть пусті

## Принцип дії скрипту
#### Наприклад у вас є 2 таблиці input.csv та dataSource.csv


###### input.csv  
id | name
------------ | -------------
1 | Cruella
2 | Spiral

###### dataSource.csv  
id | imdb
------------ | -------------
1 | 3228774
3 | 123

#### Ви отримаєте 3 таблиці

###### mappedOutputFileName.csv 
id | name | imdb
------------ | ------------- | -------------
1 | Cruella | tt3228774

###### notMappedOutputFileName.csv 
id | name | imdb
------------ | ------------- | -------------
2 | Spiral | 

###### combinedOutputFileName.csv 
id | name | imdb
------------ | ------------- | -------------
1 | Cruella | tt3228774
2 | Spiral | 

## Чому?
Для вивчення go, вирішив зробити першу програму. Ідей особливо не було, тож згадав про фільми. Сам я уже з 2015 не використовую "кинопоиск" та успішно перейшов на imdb та частково використовую tmdb.





## Просунуте використання
#### Запуск скрипту з перепризначенням файлу конфігурації
```
data-transformer.exe -config="config.json"
```

Є можливість перевизначити всі змінні(крім _convertToFullImdbString_), через запуск скрипту з додатковими параметрами, наприклад
```
data-transformer.exe -inputFileName="my_kpvotes.csv"
```


## Додаткові можливості
У файлі config.json є багато значень:
- inputFileName - ваша таблиця csv. Потрібно щоб було наявний ідентифікатор kp
- dataSourceFileName - таблиця kp2imdb.csv. Цей файл є в комплекті. Тут має бути ідентифікатор kp та imdb
- mappedOutputFileName - таблиця тільки тих ваших оцінок, де скрипт знайшов ідентифікатор imdb
- notMappedOutputFileName - таблиця тільки тих ваших оцінок, де скрипт не знайшов ідентифікатор imdb
- combinedOutputFileName - таблиця всіх ваших оцінок. Не знайдені imdb будуть пусті
- inputColumn - колонка вашого ідентифікатора. Він має називатись так само як і в kp2imdb.csv
- outputColumn - колонка вашого ідентифікатора 2, який скрипт візьме з kp2imdb.csv і додасть у вашу таблицю
- convertToFullImdbString - чи конвертувати число imdb в повний imdb. Наприклад 123 to tt0000123