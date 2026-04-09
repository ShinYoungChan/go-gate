import 'package:flutter/material.dart';
import 'qr_screen.dart';
import 'my_page_screen.dart';
import '../services/api_service.dart'; // 💡 ApiService 임포트 확인

class MainScreen extends StatefulWidget {
  const MainScreen({super.key});

  @override
  State<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
  List<dynamic> _locations = []; // 서버에서 받아올 데이터를 담을 공간
  String _userName = "로딩 중..."; // 유저 이름 초기값
  String _userEmail = "";
  String _joinDate = "";
  bool _isLoading = true;
  final String _fixedUserId = "1"; // 💡 userID 고정값 설정
  int _selectedIndex = 0;

  @override
  void initState() {
    super.initState();
    _loadInitialData(); // 💡 장소와 유저 정보를 모두 로드하도록 변경
  }

  // 💡 데이터 로드를 총괄하는 함수 (기존 _fetchLocations의 역할을 포함)
  Future<void> _loadInitialData() async {
    setState(() => _isLoading = true);

    try {
      // API 두 개를 동시에 호출해서 시간을 단축합니다.
      await Future.wait([_fetchUserInfo(), _fetchLocations()]);
    } catch (e) {
      print("전체 데이터 로드 에러: $e");
      _showErrorSnackBar("데이터를 불러오는 중 오류가 발생했습니다.");
    } finally {
      setState(() => _isLoading = false);
    }
  }

  // 1. 유저 정보 로드 (고정 ID 사용)
  Future<void> _fetchUserInfo() async {
    try {
      // ApiService에 getUserProfile(int id) 함수가 구현되어 있어야 합니다.
      final response = await ApiService().getUserInfo(
        _fixedUserId,
      ); // 💡 ID 1번 고정
      if (response != null && response.statusCode == 200) {
        setState(() {
          _userName = response.data['data']['name'] ?? "사용자";
          _userEmail = response.data['data']['email'] ?? "이메일";
          _joinDate = response.data['data']['joindate'] ?? "";
        });
      }
    } catch (e) {
      print("유저 정보 연동 에러: $e");
    }
  }

  // 2. 장소 리스트 로드 (기존 로직 유지)
  Future<void> _fetchLocations() async {
    try {
      final response = await ApiService().getLocations();
      if (response != null && response.statusCode == 200) {
        setState(() {
          _locations = response.data['data'] ?? [];
        });
      } else {
        _showErrorSnackBar("장소 정보를 가져오지 못했습니다.");
      }
    } catch (e) {
      print("장소 연동 에러: $e");
      rethrow; // _loadInitialData에서 캐치하도록 던짐
    }
  }

  void _showErrorSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: Colors.red),
    );
  }

  @override
  Widget build(BuildContext context) {
    // 탭별로 보여줄 화면 설정
    final List<Widget> _pages = [
      // [0번: 홈 탭] - 인사말을 리스트 상단에 포함
      _isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: _fetchLocations,
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  // 💡 홈 화면에서만 보이는 인사말 (AppBar에서 이사 옴)
                  Padding(
                    padding: const EdgeInsets.only(top: 20, bottom: 20),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          "안녕하세요, $_userName님!",
                          style: const TextStyle(
                            color: Colors.black,
                            fontSize: 24, // 크기를 조금 더 키워도 예쁩니다
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 6),
                        const Text(
                          "입장하실 장소를 선택하세요",
                          style: TextStyle(color: Colors.grey, fontSize: 15),
                        ),
                      ],
                    ),
                  ),

                  // 장소 리스트 렌더링
                  if (_locations.isEmpty)
                    const Center(
                      child: Padding(
                        padding: EdgeInsets.only(top: 100),
                        child: Text("등록된 장소가 없습니다."),
                      ),
                    )
                  else
                    // .map().toList()를 사용해 기존 리스트 카드들을 뿌려줍니다.
                    ..._locations
                        .map((loc) => _buildLocationCard(loc))
                        .toList(),
                ],
              ),
            ),

      // [1번: 결제 탭]
      const Center(child: Text("결제 내역")),

      // [2번: 내 정보 탭]
      MyPageScreen(
        userName: _userName,
        userEmail: _userEmail,
        joinData: _joinDate,
      ),
    ];

    return Scaffold(
      backgroundColor: Colors.white,
      // AppBar를 아예 삭제하면 본문이 폰 맨 위(시계 있는 곳)까지 올라가므로 SafeArea를 사용
      body: SafeArea(child: _pages[_selectedIndex]),

      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _selectedIndex,
        onTap: (index) {
          setState(() {
            _selectedIndex = index;
          });
        },
        backgroundColor: Colors.white,
        selectedItemColor: Colors.blue,
        unselectedItemColor: Colors.grey,
        showSelectedLabels: true,
        showUnselectedLabels: true,
        type: BottomNavigationBarType.fixed,
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home), label: "홈"),
          BottomNavigationBarItem(icon: Icon(Icons.payment), label: "결제"),
          BottomNavigationBarItem(icon: Icon(Icons.person), label: "내 정보"),
        ],
      ),
    );
  }

  Widget _buildLocationCard(dynamic loc) {
    // 💡 주의: Go에서 JSON 태그를 설정하지 않았다면 'PlaceName' (대문자)
    // 만약 `json:"place_name"` 처럼 태그를 달았다면 소문자로 써야 합니다.
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(15)),
      child: InkWell(
        onTap: () {
          // 클릭 시 해당 장소의 ID를 가지고 회원권 체크 로직으로 이동
          print("${loc['PlaceName']} (ID: ${loc['ID']}) 클릭됨");

          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => QRScreen(
                location: loc, // 클릭한 장소 정보 전체 전달
                userId: _fixedUserId, // 고정된 유저 ID (1) 전달
              ),
            ),
          );
        },
        borderRadius: BorderRadius.circular(15),
        child: Row(
          children: [
            ClipRRect(
              borderRadius: const BorderRadius.horizontal(
                left: Radius.circular(15),
              ),
              child: Image.network(
                loc['ImageURL'] ??
                    "https://via.placeholder.com/100", // 이미지 없을 시 기본값
                width: 100,
                height: 100,
                fit: BoxFit.cover,
                errorBuilder: (context, error, stackTrace) => Container(
                  width: 100,
                  height: 100,
                  color: Colors.grey[200],
                  child: const Icon(Icons.image_not_supported),
                ),
              ),
            ),
            Expanded(
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 8,
                        vertical: 4,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.blue[50],
                        borderRadius: BorderRadius.circular(5),
                      ),
                      child: Text(
                        loc['Category'] ?? "기타",
                        style: const TextStyle(
                          color: Colors.blue,
                          fontSize: 10,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      loc['PlaceName'] ?? "이름 없음",
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        const Icon(
                          Icons.location_on,
                          size: 14,
                          color: Colors.grey,
                        ),
                        const SizedBox(width: 4),
                        Expanded(
                          child: Text(
                            loc['Address'] ?? "주소 정보 없음",
                            style: const TextStyle(
                              fontSize: 12,
                              color: Colors.grey,
                            ),
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
