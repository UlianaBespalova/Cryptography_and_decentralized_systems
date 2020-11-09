# ДЗ №3

## Задание

Разработать консольное приложение, выполняющее разделение секретной строки в hex-формате длиной 256 bit (например, приватный ключ ECDSA secp256k1) на N частей по схеме Шамира и восстанавливает его при предъявлении любых T частей.

### Аргументы командной строки:
Единственный аргумент указывающий режим работы программы:

    ./Shamir split - режим разделения
    ./Shamir recover - режим восстановления

Входные данные - в режиме split (stdin):
- 1 строка - Два числа N, T, где 2 < T <= N < 100
- 2 строка - Приватный ключ (P_KEY)

Выходные данные:
- N строк, в каждой строке содержится кусочек разделенного секрета(в виде string)


Входные данные - в режиме recover (stdin):
- T или более строк с кусочками секрета (в каждой строке секрет в том же формате, что и вывод программы в режиме split)

Выходные данные:
- Приватный ключ (P_KEY), в том же виде, что и перед разделением. 


## Установка и запуск

Сборка

```
cmake CMakeLists.txt
make
```

Запуск

```
./Shamir
```

## Пример выполнения программы

Результаты разделения одного и того же секрета в разные моменты запуска приложения:
  
1)  
![example1](screenshots/split(2).png)  
-----------------------------------------------------------   
  
2)  
![example2](screenshots/split(1).png)  
  

  
-----------------------------------------------------------  
 
  
Восстановление исходного секрета по сгенерированным ключам:
  
T = 3  
![example3](screenshots/recover(1).png)  
-----------------------------------------------------------   
  
T = 4  
![example4](screenshots/recover(2).png)  
