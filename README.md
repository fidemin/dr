# Drama and company

# 설치

## 초기 설치

1. 다운 받기
```console
git clone https://github.com/yhmin84/drama.git drama
```

2. database 생성 필요 (mariadb 10.3 사용 권장. db 이름은 상관 없다.)
```
CREATE DATABASE `dramadb`;
```

3. configure 파일을 생성한다. json 파일이고, 생성 위치는 상관없다.

	```javascript
	{
		"port": 5001,
		"debug": true,
		"secret": "R3sCZY2QmuvMRg4U8XLLbvwWGHKrDfko",
		"db":{
			"host": "127.0.0.1",
			"port": 3306,
			"name": "dramadb",
			"user": "drama",
			"pass": "dramapass"
		}
	}
	```

	**config 파일의 각 요소 설명**

	`port`: APP 이 실행되는 포트 

	`debug`: 디버깅 레벨을 결정한다. 프로덕션에서는 `false` 사용을 권장한다.

	`secret`: JWT 시그니처 만들때 사용하는 키

	`db`: mysql DB 연결 정보를 넣어준다. 


4. database 초기화 및 데모 데이터 넣기 

	* macos 에서

	```
	./bin/macos/initdb --config=config_file_path
	```

	* linux 에서
	
	```
	./bin/linux/initdb.linux.amd64 --config=config_file_path
	```


## 실행 하기
	* macos 에서

	```
	./bin/macos/start --config=config_file_path
	```

	* linux 에서

	```
	./bin/linux/start.linux.amd64 --config=config_file_path
	```


## API 문서
* [링크](https://documenter.getpostman.com/view/2460249/S17jWXaS)

