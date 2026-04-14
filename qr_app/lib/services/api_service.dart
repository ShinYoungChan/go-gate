import 'package:dio/dio.dart';

class ApiService {
  // 어디서든 'ApiService()'로 불러서 쓸 수 있게 만드는 설정입니다.
  static final ApiService _instance = ApiService._internal();
  factory ApiService() => _instance;

  late Dio dio;

  ApiService._internal() {
    dio = Dio(
      BaseOptions(
        // 안드로이드 에뮬레이터라면 10.0.2.2, iOS라면 localhost를 사용
        baseUrl: 'http://localhost:8080',
        contentType: 'application/json',
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
        data: {'email': email, 'password': password},
      );
      return response;
    } on DioException catch (e) {
      // 에러가 나면 콘솔에 찍어줍니다.
      print("통신 에러 발생: ${e.message}");
      return e.response;
    }
  }

  Future<Response?> getUserInfo(String userId) async {
    try {
      // URL 뒤에 /:id 대신 실제 ID를 붙여서 요청합니다.
      // 예: /user/info/123
      return await dio.get('/user/info/$userId');
    } on DioException catch (e) {
      print("유저 정보 로드 실패 (ID: $userId): ${e.message}");
      return e.response;
    }
  }

  Future<Response?> checkUserMembership(String userId, String locationId) async {
    try {
      return await dio.get('/membership/info/$userId/$locationId');
    } on DioException catch (e) {
      print("유저 멤버십 정보 로드 실패 (ID: $userId), (LocationID: $locationId): ${e.message}");
      return e.response;
    }
  }

  Future<Response?> getQrCode(String userId, String locationId) async {
    try {
      return await dio.post(
        '/api/v1/entry/token',
        data: {"user_id": int.parse(userId), "location_id": int.parse(locationId)},
      );
    } on DioException catch (e) {
      print("QR 정보 로드 실패: ${e.message}");
      return e.response;
    }
  }

  Future<Response?> getLocations() async{
    try{
      return await dio.get('/api/v1/locations');
    }on DioException catch(e){
      print("장소 정보 로드 실패: ${e.message}");
      return e.response;
    }
  }

  Future<Response?> getUserSummary(String userId) async{
    try{
      return await dio.get('/user/mypage/$userId');
    }on DioException catch(e){
      print("유저 정보 로드 실패: ${e.message}");
      return e.response;
    }
  }

  Future<Response?> getLocationMemberships(String locationId) async{
    try{
      return await dio.get('/location/membership/$locationId');
    }on DioException catch(e){
      print("멤버쉽 정보 로드 실패: ${e.message}");
      return e.response;
    }
  }
}
