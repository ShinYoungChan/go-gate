import 'package:flutter/material.dart';
import 'screens/login_screen.dart'; // 방금 만든 화면 불러오기

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'QR App',
      theme: ThemeData(primarySwatch: Colors.blue),
      // 앱이 시작되면 LoginScreen을 가장 먼저 보여줍니다!
      home: const LoginScreen(), 
    );
  }
}