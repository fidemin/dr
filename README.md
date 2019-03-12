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


4. database 테이블 생성 

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

## 테스트
* 테스트는 golang이 설치되어야 실행이 가능하기 때문에, 결과 문서로 대체한다.
  * [테스트 결과 문서](./TESTRESULT.txt) 

* 테스트에는 db mock 객체를 이용한 유닛테스트만 작성하였다. 실제 db와 연동해서 진행하는 integration 테스트는 작성하지 않았다. `_test.go`로 끝나는 파일이 테스트 파일이다.

* go가 설치되어 있다면, 다음과 같은 명령어로 테스트 할 수 있다. (glide라는 go 패키지 관리 툴을 설치해야 한다.)

```console
brew install glide
cd $GOPATH/src
go get -u github.com/yhmin84/drama
cd drama
glide install
go test
```

## 디렉토리 및 파일 설명
### golang 파일
* `main.go` : 실행 파일
* `auth.handler.go`: 로그인/회원가입 API 코드가 정의되어 있는 파일
* `auth.interface.go`: 유저 객체에 대한 비즈니스 로직 인터페이스와 이를 구현한 DB 객체가 있음
* `auth.model.go`: 유저 ORM 객체와 타입 객체 정의가 들어있다.
* `auth_test.go`: 유저 API 테스트 파일. mock 객체도 구현되어 있다.
* `dispatch.handler.go`: 배차 관련 API 코드가 정의되어 있는 파일
* `dispatch.interface.go`: 배차 객체에 대한 비즈니스 로직 인터페이스와 이를 구현한 DB 객체가 있음
* `dispatch.model.go` : 배차 ORM 객체와 타입 객체 정의가 들어 있다.
* `dispatch_test.go` : 배차 API 테스트 파일. mock 객체도 구현되어 있다.
* `jwt.go`: jwt 타입 정의와 함수 구현
* `util.go`: 유틸 파일
* `util_test.go`: 유틸 테스트 파일

### 기타
* `glide.yaml` : javascript의 package.json과 비슷한 기능을 하는 파일
* `glide.lock` : javascript의 package.lock과 비슷한 기능을 하는 파일
* `bin` : macos와 linux 환경에서 작동하는 `initdb`, `start` 실행 파일이 들어있음
* `sh` : 개발시 사용하는 유용한 bash script 파일이 들어있음
* `init` : db 테이블 생성 코드가 들어있음





