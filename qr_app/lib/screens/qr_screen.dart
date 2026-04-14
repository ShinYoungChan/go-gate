import 'package:flutter/material.dart';
import 'package:qr_flutter/qr_flutter.dart'; // 패키지 추가 필요
import '../services/api_service.dart';
import 'membership_purchase_screen.dart';
import 'dart:async';

class QRScreen extends StatefulWidget {
  final dynamic location;
  final String userId;

  const QRScreen({super.key, required this.location, required this.userId});

  @override
  State<QRScreen> createState() => _QRScreenState();
}

class _QRScreenState extends State<QRScreen> {
  bool _hasPass = false;
  bool _isLoading = true;
  String _qrData = ""; // API에서 받은 토큰 저장
  int _secondsRemaining = 30;
  Timer? _countdownTimer;
  Map<String, dynamic>? _membershipData;

  @override
  void initState() {
    super.initState();
    _checkAndInitialize();
  }

  @override
  void dispose() {
    _countdownTimer?.cancel(); // 화면 나갈 때 타이머 해제 필수
    super.dispose();
  }

  // 처음 진입 시 이용권 확인 및 QR 초기 발급
  Future<void> _checkAndInitialize() async {
    setState(() => _isLoading = true);
    try {
      // 1. 이용권 확인
      final passRes = await ApiService().checkUserMembership(
        widget.userId,
        widget.location['ID'].toString(),
      );

      if (passRes != null && passRes.statusCode == 200) {
        bool isValid = passRes.data['data']['IsValid'] ?? false;
        setState(() {
          _membershipData = passRes.data['data'];

          _hasPass = _membershipData!['IsValid'] ?? false;
        });

        print(passRes.data['data']);

        // 2. 이용권이 있다면 바로 QR 발급
        if (isValid) {
          await _refreshQrOnly();
        }
      }
    } catch (e) {
      print("초기화 실패: $e");
    } finally {
      setState(() => _isLoading = false);
    }
  }

  // 💡 작성해주신 기능: 30초 타이머 로직
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

  // 💡 작성해주신 기능: QR 코드 단독 갱신 (서버 API 호출)
  Future<void> _refreshQrOnly() async {
    try {
      final response = await ApiService().getQrCode(
        widget.userId,
        widget.location['ID'].toString(),
      );

      if (response != null && response.statusCode == 200) {
        setState(() {
          // 서버에서 받은 실제 토큰값으로 QR 생성
          _qrData = response.data['data']['token'];
          _startTimer(); // 갱신 후 타이머 재시작
        });
      }
    } catch (e) {
      print("QR 갱신 실패: $e");
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text(
            "< 장소 선택으로 돌아가기",
            style: TextStyle(color: Colors.black, fontSize: 12),
          ),
        ),
        leadingWidth: 200,
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              child: Column(
                children: [
                  _buildHeaderCard(),
                  const SizedBox(height: 20),
                  // 이용권 상태에 따른 화면 분기
                  _hasPass ? _buildQRSection() : _buildNoPassSection(),
                  const SizedBox(height: 40),
                  _buildFooterInstructions(),
                ],
              ),
            ),
    );
  }

  // 상단 장소 요약 (사진 상단 부분)
  Widget _buildHeaderCard() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(15),
          border: Border.all(color: Colors.grey[200]!),
        ),
        child: Row(
          children: [
            ClipRRect(
              borderRadius: BorderRadius.circular(10),
              child: Image.network(
                widget.location['ImageURL'],
                width: 60,
                height: 60,
                fit: BoxFit.cover,
              ),
            ),
            const SizedBox(width: 15),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  widget.location['PlaceName'],
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Text(
                  widget.location['Address'],
                  style: const TextStyle(fontSize: 12, color: Colors.grey),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  // 수정된 QR 섹션 (만료시간 + 멤버십 상세 정보 포함)
  Widget _buildQRSection() {
    // 데이터가 없을 경우를 대비한 방어 코드
    if (_membershipData == null) return const SizedBox();
    // 멤버십 데이터 변수
    final membership = _membershipData!; // 서버에서 받은 data 객체
    final bool isCountPass = membership['IsCountType']; // 서버 필드명에 맞게 수정

    return Column(
      children: [
        const Icon(Icons.qr_code_2, size: 30),
        const Text(
          "입장용 QR 코드",
          style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const Text(
          "입구에서 이 QR 코드를 스캔해주세요",
          style: TextStyle(color: Colors.grey),
        ),
        const SizedBox(height: 20),

        // QR 코드 본체
        _qrData.isNotEmpty
            ? QrImageView(data: _qrData, size: 200)
            : const SizedBox(
                height: 200,
                child: Center(child: CircularProgressIndicator()),
              ),

        const SizedBox(height: 10),

        // 1. 만료 시간 (실시간 타이머 반영)
        Text(
          "남은 시간: $_secondsRemaining초",
          style: const TextStyle(
            color: Colors.red,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 20),

        // 2. 멤버십 상세 정보 카드 (요청하신 부분)
        Container(
          padding: const EdgeInsets.all(16),
          margin: const EdgeInsets.symmetric(horizontal: 20),
          decoration: BoxDecoration(
            color: Colors.grey[100],
            borderRadius: BorderRadius.circular(12),
          ),
          child: Column(
            children: [
              // 횟수권/정기권 구분 태그
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: Colors.blue,
                  borderRadius: BorderRadius.circular(4),
                ),
                child: Text(
                  isCountPass ? "횟수권" : "정기권",
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 12,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              const SizedBox(height: 10),

              // 횟수권일 때만 잔여 횟수 표기
              if (isCountPass) ...[
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Text("잔여 횟수: ", style: TextStyle(fontSize: 16)),
                    Text(
                      "${membership['Count']}회",
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                        color: Colors.blue,
                      ),
                    ),
                  ],
                ),
                const Divider(),
              ],

              // 사용 기한 표기
              Text(
                "사용기한: ${membership['SttDt'].split('T')[0]} ~ ${membership['EndDt'].split('T')[0]}",
                style: const TextStyle(fontSize: 13, color: Colors.black87),
              ),
            ],
          ),
        ),

        const SizedBox(height: 10),
        TextButton.icon(
          onPressed: _refreshQrOnly,
          icon: const Icon(Icons.refresh),
          label: const Text("QR 코드 즉시 갱신"),
        ),
      ],
    );
  }

  // 이용권이 없을 때 보여줄 구매 영역 (기존 코드 유지)
  Widget _buildNoPassSection() {
    return Column(
      children: [
        const Icon(Icons.error_outline, size: 60, color: Colors.orange),
        const SizedBox(height: 16),
        const Text(
          "유효한 이용권이 없습니다.",
          style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 24),
        ElevatedButton(
          onPressed: () {
            Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) => MembershipPurchaseScreen(
                  locationId: widget.location['ID'], // 현재 QR 화면이 들고 있는 지점 ID 전달
                  locationName: widget.location['PlaceName'], // (선택) 화면 상단에 보여줄 이름
                ),
              ),
            );
          },
          style: ElevatedButton.styleFrom(
            backgroundColor: Colors.blue,
            minimumSize: const Size(200, 50),
          ),
          child: const Text("이용권 구매하기", style: TextStyle(color: Colors.white)),
        ),
      ],
    );
  }

  // 하단 이용안내 문구
  Widget _buildFooterInstructions() {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(20),
      color: Colors.grey[50],
      child: const Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text("이용 안내", style: TextStyle(fontWeight: FontWeight.bold)),
          SizedBox(height: 10),
          Text(
            "• 해당 QR 코드는 1회성 입장용입니다.",
            style: TextStyle(fontSize: 12, color: Colors.grey),
          ),
          Text(
            "• 캡쳐된 화면으로는 입장이 불가할 수 있습니다.",
            style: TextStyle(fontSize: 12, color: Colors.grey),
          ),
        ],
      ),
    );
  }
}
