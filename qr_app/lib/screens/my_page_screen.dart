import 'package:flutter/material.dart';
import '../services/api_service.dart';
import 'package:intl/intl.dart';

class MyPageScreen extends StatefulWidget {
  final String userName;
  final String userEmail;
  final String joinData;

  const MyPageScreen({
    super.key,
    required this.userName,
    required this.userEmail,
    required this.joinData,
  });

  @override
  State<MyPageScreen> createState() => _MyPageScreenState();
}

class _MyPageScreenState extends State<MyPageScreen> {
  int _entryCount = 0;
  int _totalAmount = 0;
  bool _isLoading = true;
  final String _fixedUserId = "1"; // 💡 userID 고정값 설정

  @override
  void initState() {
    super.initState();
    _fetchMyPageSummary();
  }

  Future<void> _fetchMyPageSummary() async {
    setState(() => _isLoading = true);

    try {
      final response = await ApiService().getUserSummary(_fixedUserId);

      if (response != null && response.statusCode == 200) {
        setState(() {
          final data = response.data?['data'];
          if (data != null) {
            _entryCount = data['entry_count'] ?? 0;
            _totalAmount = data['total_amount'] ?? 0;
          }
        });
      }
    } catch (e) {
      print("데이터 로드 에러: $e");
      _showErrorSnackBar("데이터를 불러오는 중 오류가 발생했습니다.");
    } finally {
      setState(() => _isLoading = false);
    }
  }

  void _showErrorSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: Colors.red),
    );
  }

  String formatJoinDate(String rawDate) {
    // rawDate가 "2026-01-01 01:00:00 ..." 일 때
    List<String> parts = rawDate
        .split(' ')[0]
        .split('-'); // ['2026', '01', '01']
    return "${parts[0]}년 ${parts[1]}월 ${parts[2]}일";
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(
        body: Center(child: CircularProgressIndicator(color: Colors.blue)),
      );
    }
    return SingleChildScrollView(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 1. 헤더 영역
          const Text(
            "내 정보",
            style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 4),
          const Text(
            "계정 정보 및 이용 현황을 확인하세요",
            style: TextStyle(color: Colors.grey, fontSize: 16),
          ),
          const SizedBox(height: 30),

          // 2. 프로필 카드 영역
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(20),
              border: Border.all(color: Colors.grey.shade200),
            ),
            child: Column(
              children: [
                Row(
                  children: [
                    CircleAvatar(
                      radius: 35,
                      backgroundColor: Colors.blue.shade600,
                      child: const Text(
                        "1",
                        style: TextStyle(color: Colors.white, fontSize: 24),
                      ),
                    ),
                    const SizedBox(width: 20),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            Text(
                              widget.userName,
                              style: TextStyle(
                                fontSize: 20,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                            const SizedBox(width: 8),
                            Container(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 8,
                                vertical: 2,
                              ),
                              decoration: BoxDecoration(
                                color: Colors.grey.shade200,
                                borderRadius: BorderRadius.circular(10),
                              ),
                              child: const Text(
                                "정회원",
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.grey,
                                ),
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 10),
                        Row(
                          children: [
                            const Icon(
                              Icons.email_outlined,
                              size: 16,
                              color: Colors.grey,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              widget.userEmail,
                              style: TextStyle(color: Colors.grey),
                            ),
                          ],
                        ),
                        const SizedBox(height: 4),
                        Row(
                          children: [
                            const Icon(
                              Icons.calendar_today_outlined,
                              size: 16,
                              color: Colors.grey,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              "가입일: ${formatJoinDate(widget.joinData)}",
                              style: TextStyle(color: Colors.grey),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(height: 20),

          // 3. 통계 카드 영역 (이번 달 입장 / 총 결제액)
          Row(
            children: [
              _buildStatCard("이번 달 입장", "$_entryCount회", Icons.qr_code_scanner),
              const SizedBox(width: 15),
              _buildStatCard("총 결제액", "${NumberFormat('#,###').format(_totalAmount)}원", Icons.credit_card),
            ],
          ),
          const SizedBox(height: 30),

          // 4. 계정 설정 영역
          const Text(
            "계정 설정",
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 15),
          Container(
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(15),
              border: Border.all(color: Colors.grey.shade200),
            ),
            child: Column(
              children: [
                _buildSettingsItem(Icons.person_outline, "프로필 수정"),
                const Divider(height: 1),
                _buildSettingsItem(Icons.payment_outlined, "결제 수단 관리"),
              ],
            ),
          ),
          const SizedBox(height: 30),

          // 5. 로그아웃 버튼
          SizedBox(
            width: double.infinity,
            height: 55,
            child: ElevatedButton.icon(
              onPressed: () {
                // 로그아웃 로직
              },
              icon: const Icon(Icons.logout, color: Colors.white),
              label: const Text(
                "로그아웃",
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFD32F2F),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
            ),
          ),
          const SizedBox(height: 50),
        ],
      ),
    );
  }

  // 통계 카드 위젯 헬퍼
  Widget _buildStatCard(String title, String value, IconData icon) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(20),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(20),
          border: Border.all(color: Colors.grey.shade200),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(icon, size: 18, color: Colors.grey),
                const SizedBox(width: 6),
                Text(
                  title,
                  style: const TextStyle(color: Colors.grey, fontSize: 14),
                ),
              ],
            ),
            const SizedBox(height: 15),
            Text(
              value,
              style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
            ),
          ],
        ),
      ),
    );
  }

  // 설정 아이템 위젯 헬퍼
  Widget _buildSettingsItem(IconData icon, String title) {
    return ListTile(
      leading: Icon(icon, color: Colors.black87),
      title: Text(title, style: const TextStyle(fontSize: 16)),
      trailing: const Icon(Icons.chevron_right, color: Colors.grey),
      onTap: () {},
    );
  }
}
