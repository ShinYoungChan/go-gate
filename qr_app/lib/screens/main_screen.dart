import 'package:flutter/material.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'dart:async';
import '../services/api_service.dart'; // 기존에 만든 서비스 파일 연결

class MainScreen extends StatefulWidget {
  const MainScreen({super.key});

  @override
  State<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
  // --- [데이터 영역] ---
  String _userName = "로드 중...";
  String _qrData = "";
  bool _hasMembership = false;
  bool _isLoading = true;

  // 멤버십 상세 정보
  bool _isValid = false; // IsValid (회원권 유효 여부)
  bool _isCountType = false; // IsCountType (true: 횟수권, false: 정기권)
  String _membershipName =
      "TestName"; // 상품명 (이건 ItemID로 조회한 상품, 테스트를 위해 Test 고정)
  int _count = 0; // Count (잔여 횟수)
  String _expiryDate = ""; // EndDt (만료 날짜)

  // --- [타이머 영역: 30초 갱신] ---
  Timer? _countdownTimer;
  int _secondsRemaining = 30; // 💡 요구사항 #1: 30초 주기

  @override
  void initState() {
    super.initState();
    _fetchAllData(); // 시작하자마자 3개 API 호출
  }

  @override
  void dispose() {
    _countdownTimer?.cancel();
    super.dispose();
  }

  /// **기능: 3개의 API(User, Membership, Entry) 동시 호출**
  Future<void> _fetchAllData() async {
    // 테스트용 고정 ID
    const String testId = "1";
    const String testLocationId = "6";
    setState(() => _isLoading = true);
    try {
      // 💡 기술 포인트: Future.wait로 3번의 통신을 병렬 처리
      final results = await Future.wait([
        ApiService().getUserInfo(testId), // [0] 유저 이름
        ApiService().getMembershipInfo(testId), // [1] 회원권 상태/상세
        ApiService().getQrCode(testId,testLocationId), // [2] QR 데이터
      ]);

      setState(() {
        // (1) 유저 정보
        if (results[0] != null && results[0]!.statusCode == 200) {
          _userName = results[0]!.data['data']['name'] ?? "회원";
        }

        // (2) 멤버십 정보
        if (results[1] != null && results[1]!.statusCode == 200) {
          var mData = results[1]!.data['data'];
          _isValid = mData['IsValid'] ?? false;
          _isCountType = mData['IsCountType'] ?? false;
          _count = mData['Count'] ?? 0;
          String rawDate = mData['EndDt'] ?? "-";
          _expiryDate = rawDate.length >= 10
              ? rawDate.substring(0, 10)
              : rawDate;

          _hasMembership = _isValid;
        }

        // (3) QR 정보
        if (results[2] != null && results[2]!.statusCode == 200) {
          var  qData = results[2]!.data['data'];
          _qrData = qData['token'] ?? "";
        }

        _isLoading = false;
        if (_hasMembership) _startTimer(); // 회원권 있을 때만 타이머 시작
      });
    } catch (e) {
      debugPrint("데이터 로드 실패: $e");
      setState(() => _isLoading = false);
    }
  }

  /// **기능: 30초 타이머 로직**
  void _startTimer() {
    _countdownTimer?.cancel();
    _secondsRemaining = 30;
    _countdownTimer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (_secondsRemaining > 0) {
        setState(() => _secondsRemaining--);
      } else {
        _refreshQrOnly(); // 0초 되면 QR만 다시 가져옴
      }
    });
  }

  /// **기능: QR 코드만 단독 갱신**
  Future<void> _refreshQrOnly() async {
    // 테스트용 고정 ID
    const String testId = "1";
    const String testLocationId = "6";
    final response = await ApiService().getQrCode(testId,testLocationId);
    if (response != null && response.statusCode == 200) {
      setState(() {
        _qrData = response.data['data']['token'];
        _startTimer();
      });
    }
  }

  // --- [UI 그리기] ---

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    return Scaffold(
      backgroundColor: const Color(0xFFF7F7FB),
      bottomNavigationBar: _buildBottomNavBar(),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(horizontal: 24.0),
          child: Column(
            children: [
              const SizedBox(height: 40),
              _buildHeader(),
              const SizedBox(height: 30),

              // 💡 요구사항 #2: 회원권 유무에 따른 분기 처리
              _hasMembership
                  ? _buildActiveMembershipUI()
                  : _buildNoMembershipUI(),

              const SizedBox(height: 20),
              _buildUsageInfoCard(),
              const SizedBox(height: 30),
            ],
          ),
        ),
      ),
    );
  }

  // (1) 헤더: 안녕하세요 님!
  Widget _buildHeader() {
    return Center(
      child: Column(
        children: [
          Text(
            "안녕하세요, $_userName님!",
            style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 4),
          const Text(
            "QR 코드로 빠르게 입장하세요",
            style: TextStyle(color: Colors.black54),
          ),
        ],
      ),
    );
  }

  // (2) 회원권이 있을 때: QR 카드 + 멤버십 정보
  Widget _buildActiveMembershipUI() {
    return Column(
      children: [
        _buildQrCard(), // 사진과 동일한 QR 카드
        const SizedBox(height: 16),
        _buildMembershipDetailCard(), // 💡 새로 추가된 멤버십 상세 카드
      ],
    );
  }

  // (3) 회원권이 없을 때: 구매 버튼만 노출
  Widget _buildNoMembershipUI() {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(30),
      decoration: _buildBoxDecoration(),
      child: Column(
        children: [
          Icon(Icons.info_outline, size: 50, color: Colors.grey[400]),
          const SizedBox(height: 16),
          const Text(
            "유효한 회원권이 없습니다.",
            style: TextStyle(fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 20),
          ElevatedButton(
            onPressed: () {}, // 회원권 구매 페이지 이동 로직
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFF1F4EF5),
              minimumSize: const Size(double.infinity, 55),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
            child: const Text(
              "회원권 구매하기",
              style: TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        ],
      ),
    );
  }

  // (4) 사진 속 디자인: 입장용 QR 카드
  Widget _buildQrCard() {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: _buildBoxDecoration(),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: const [
              Icon(Icons.qr_code_2, size: 20),
              SizedBox(width: 8),
              Text("입장용 QR 코드", style: TextStyle(fontWeight: FontWeight.bold)),
            ],
          ),
          const Text(
            "입구에서 이 QR 코드를 스캔해주세요",
            style: TextStyle(fontSize: 12, color: Colors.black54),
          ),
          const SizedBox(height: 24),
          Center(child: QrImageView(data: _qrData, size: 180)),
          const SizedBox(height: 16),
          Center(
            child: Column(
              children: [
                Text(
                  "남은 시간: $_secondsRemaining초",
                  style: TextStyle(
                    fontSize: 13,
                    color: Colors.blue[700],
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                OutlinedButton.icon(
                  onPressed: _refreshQrOnly,
                  icon: const Icon(Icons.refresh, size: 16),
                  label: const Text("QR 코드 갱신"),
                  style: OutlinedButton.styleFrom(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(20),
                    ),
                    foregroundColor: Colors.black87,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  // (5) 💡 멤버십 상세 카드 (요구사항 #3: 종류, 잔여횟수, 기간)
  Widget _buildMembershipDetailCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: _buildBoxDecoration(),
      child: Row(
        children: [
          // 아이콘 (횟수권은 숫자 아이콘, 정기권은 달력 아이콘으로 차별화 가능)
          Container(
            padding: const EdgeInsets.all(10),
            decoration: BoxDecoration(
              color: Colors.blue[50],
              shape: BoxShape.circle,
            ),
            child: Icon(
              _isCountType ? Icons.confirmation_number : Icons.calendar_today,
              color: const Color(0xFF1F4EF5),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // 1. 회원권 종류 표기
                Text(
                  _isCountType ? "횟수권 이용 중" : "정기권 이용 중",
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 15,
                  ),
                ),
                const SizedBox(height: 4),

                // 2. 조건부 상세 정보 표기 (💡 핵심 로직)
                Text(
                  _isCountType
                      ? "잔여 $_count회 / $_expiryDate까지" // 횟수권일 때
                      : "$_expiryDate까지", // 정기권일 때
                  style: const TextStyle(fontSize: 13, color: Colors.black54),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  // (6) 이용 안내 카드
  Widget _buildUsageInfoCard() {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(20),
      decoration: _buildBoxDecoration(),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text("이용 안내", style: TextStyle(fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),
          _buildBulletItem("QR 코드는 30초마다 자동으로 갱신됩니다."),
          _buildBulletItem("입장 시 화면을 밝게 조절해주세요."),
        ],
      ),
    );
  }

  Widget _buildBulletItem(String text) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 6),
      child: Row(
        children: [
          const Icon(Icons.circle, size: 4, color: Color(0xFF1F4EF5)),
          const SizedBox(width: 8),
          Text(
            text,
            style: const TextStyle(fontSize: 12, color: Colors.black54),
          ),
        ],
      ),
    );
  }

  BoxDecoration _buildBoxDecoration() {
    return BoxDecoration(
      color: Colors.white,
      borderRadius: BorderRadius.circular(16),
      boxShadow: [
        BoxShadow(
          color: Colors.black.withOpacity(0.05),
          blurRadius: 10,
          offset: const Offset(0, 4),
        ),
      ],
    );
  }

  Widget _buildBottomNavBar() {
    return BottomNavigationBar(
      selectedItemColor: const Color(0xFF1F4EF5),
      unselectedItemColor: Colors.grey,
      items: const [
        BottomNavigationBarItem(icon: Icon(Icons.home), label: "홈"),
        BottomNavigationBarItem(icon: Icon(Icons.payment), label: "결제"),
        BottomNavigationBarItem(icon: Icon(Icons.person), label: "내정보"),
      ],
    );
  }
}
