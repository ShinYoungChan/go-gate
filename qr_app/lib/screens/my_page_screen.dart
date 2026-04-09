import 'package:flutter/material.dart';

class MyPageScreen extends StatelessWidget {
  final String userName;
  final String userEmail;
  final String joinData;

  const MyPageScreen({
    super.key,
    required this.userName,
    required this.userEmail,
    required this.joinData,
  });

  String formatJoinDate(String rawDate) {
    // rawDate가 "2026-01-01 01:00:00 ..." 일 때
    List<String> parts = rawDate.split(' ')[0].split('-'); // ['2026', '01', '01']
    return "${parts[0]}년 ${parts[1]}월 ${parts[2]}일";
  }

  @override
  Widget build(BuildContext context) {
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
                      child: const Text("1", style: TextStyle(color: Colors.white, fontSize: 24)),
                    ),
                    const SizedBox(width: 20),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            Text("$userName", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
                            const SizedBox(width: 8),
                            Container(
                              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                              decoration: BoxDecoration(
                                color: Colors.grey.shade200,
                                borderRadius: BorderRadius.circular(10),
                              ),
                              child: const Text("정회원", style: TextStyle(fontSize: 12, color: Colors.grey)),
                            ),
                          ],
                        ),
                        const SizedBox(height: 10),
                        Row(
                          children: [
                            const Icon(Icons.email_outlined, size: 16, color: Colors.grey),
                            const SizedBox(width: 4),
                            Text("$userEmail", style: TextStyle(color: Colors.grey)),
                          ],
                        ),
                        const SizedBox(height: 4),
                        Row(
                          children: [
                            const Icon(Icons.calendar_today_outlined, size: 16, color: Colors.grey),
                            const SizedBox(width: 4),
                            Text("가입일: ${formatJoinDate(joinData)}", style: TextStyle(color: Colors.grey)),
                          ],
                        ),
                      ],
                    )
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(height: 20),

          // 3. 통계 카드 영역 (이번 달 입장 / 총 결제액)
          Row(
            children: [
              _buildStatCard("이번 달 입장", "12회", Icons.qr_code_scanner),
              const SizedBox(width: 15),
              _buildStatCard("총 결제액", "180,000원", Icons.credit_card),
            ],
          ),
          const SizedBox(height: 30),

          // 4. 계정 설정 영역
          const Text("계정 설정", style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
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
              label: const Text("로그아웃", style: TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.bold)),
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFD32F2F),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
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
                Text(title, style: const TextStyle(color: Colors.grey, fontSize: 14)),
              ],
            ),
            const SizedBox(height: 15),
            Text(value, style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold)),
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