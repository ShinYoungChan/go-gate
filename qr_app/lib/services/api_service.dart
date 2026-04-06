import 'package:dio/dio.dart';

class ApiService {
  // 어디서든 'ApiService()'로 불러서 쓸 수 있게 만드는 설정입니다.
  static final ApiService _instance = ApiService._internal();
  factory ApiService() => _instance;

  late Dio dio;

  ApiService._internal() {
    dio = Dio(
      BaseOptions(
        // ⚠️ 안드로이드 에뮬레이터라면 10.0.2.2, iOS라면 localhost를 쓰세요!
        baseUrl: 'http://localhost:8080', 
        connectTimeout: const Duration(seconds: 5),
        receiveTimeout: const Duration(seconds: 3),
      ),
    );
  }

  // 로그인 요청 함수
  Future<Response?> login(String email, String password) async {
    try {
      final response = await dio.post(
        '/login', // Go 백엔드에서 만든 로그인 경로
        data: {
          'email': email,
          'password': password,
        },
      );
      return response;
    } on DioException catch (e) {
      // 에러가 나면 콘솔에 찍어줍니다.
      print("통신 에러 발생: ${e.message}");
      return e.response;
    }
  }
}