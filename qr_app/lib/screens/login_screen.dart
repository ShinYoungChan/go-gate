import 'package:flutter/material.dart';
import 'package:dio/dio.dart'; // 아까 설치한 dio
import 'package:qr_app/services/api_service.dart';
import 'main_screen.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();

  // 로그인 버튼 클릭 시 실행될 함수
  void _handleLogin() async {
    final email = _emailController.text;
    final password = _passwordController.text;

    print("로그인 시도: $email / $password");

    // 서버에 요청 보내기
    final response = await ApiService().login(email, password);

    if (response != null && response.statusCode == 200) {
      // ✅ 로그인 성공 시 화면 이동!
      if (!mounted) return; // 위젯이 화면에 살아있는지 확인 (안전장치)

      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (context) => const MainScreen()),
      );
    } else {
      // 로그인 실패 알림
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('로그인 실패: 정보를 확인해주세요.')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('QR 로그인')),
      body: Padding(
        padding: const EdgeInsets.all(20.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            TextField(
              controller: _emailController,
              decoration: const InputDecoration(labelText: '이메일'),
            ),
            const SizedBox(height: 10),
            TextField(
              controller: _passwordController,
              decoration: const InputDecoration(labelText: '비밀번호'),
              obscureText: true, // 비밀번호 가리기
            ),
            const SizedBox(height: 30),
            ElevatedButton(
              onPressed: _handleLogin,
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 50),
              ),
              child: const Text('로그인'),
            ),
          ],
        ),
      ),
    );
  }
}
