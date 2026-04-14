import 'package:flutter/material.dart';
import 'package:intl/intl.dart'; // 금액 포맷팅용 (pubspec.yaml에 추가 필요)
import 'dart:convert';

import 'package:qr_app/services/api_service.dart';
import '../services/api_service.dart';

class MembershipPurchaseScreen extends StatefulWidget {
  final int locationId;
  final String locationName;

  const MembershipPurchaseScreen({
    super.key,
    required this.locationId,
    required this.locationName,
  });

  @override
  State<MembershipPurchaseScreen> createState() =>
      _MembershipPurchaseScreenState();
}

class _MembershipPurchaseScreenState extends State<MembershipPurchaseScreen> {
  List<dynamic> _items = [];
  bool _isLoading = true;
  int? _selectedItemId;
  dynamic _selectedItem;

  final NumberFormat _currencyFormat = NumberFormat.simpleCurrency(
    locale: 'ko_KR',
  );

  @override
  void initState() {
    super.initState();
    _fetchMembershipItems();
  }

  // 1. 백엔드에서 지점별 상품 리스트 가져오기
  Future<void> _fetchMembershipItems() async {
    try {
      // 본인의 서버 주소로 변경하세요 (예: http://localhost:8080)
      final response = await ApiService().getLocationMemberships(
        widget.locationId.toString(),
      );

      if (response != null && response.statusCode == 200) {
        setState(() {
          // 서버에서 받은 실제 토큰값으로 QR 생성
          _items = response.data['data'];
          print(_items);
        });
      }
    } catch (e) {
      setState(() => _isLoading = false);
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('오류 발생: $e')));
    } finally {
      // 💡 성공하든 실패하든 결국 로딩 상태를 false로 변경!
      setState(() {
        _isLoading = false;
      });
    }
  }

  // 2. 결제 버튼 클릭 시 실행될 함수
  void _handlePayment() {
    // 여기서 나중에 실제 결제 모듈(토스, 포트원 등)을 붙입니다.
    // 지금은 성공했다고 가정하고 이전 화면으로 돌아가는 로직만!
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text("결제 확인"),
        content: Text("${_selectedItem['Title']}을(를) 구매하시겠습니까?"),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text("취소"),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context); // 다이얼로그 닫기
              Navigator.pop(context, true); // QR 화면으로 '성공' 신호(true) 보내며 돌아가기
            },
            child: const Text("결제하기"),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("${widget.locationName} 이용권 구매"),
        centerTitle: true,
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _items.isEmpty
          ? const Center(child: Text("현재 판매 중인 이용권이 없습니다."))
          : ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _items.length,
              itemBuilder: (context, index) {
                final item = _items[index];
                final isSelected = _selectedItemId == item['ID'];

                return Card(
                  margin: const EdgeInsets.only(bottom: 12),
                  elevation: isSelected ? 4 : 1,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                    side: BorderSide(
                      color: isSelected ? Colors.blue : Colors.transparent,
                      width: 2,
                    ),
                  ),
                  child: ListTile(
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 20,
                      vertical: 8,
                    ),
                    title: Text(
                      item['Title'],
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 18,
                      ),
                    ),
                    subtitle: Text("${item['DurationDays']}일 이용 가능"),
                    trailing: Text(
                      _currencyFormat.format(item['Amount']),
                      style: const TextStyle(
                        color: Colors.blueAccent,
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                    onTap: () {
                      setState(() {
                        _selectedItemId = item['ID'];
                        _selectedItem = item;
                      });
                    },
                  ),
                );
              },
            ),
      // 3. 하단 고정 결제 버튼
      bottomNavigationBar: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: ElevatedButton(
            onPressed: _selectedItemId == null ? null : _handlePayment,
            style: ElevatedButton.styleFrom(
              minimumSize: const Size(double.infinity, 56),
              backgroundColor: Colors.blue,
              disabledBackgroundColor: Colors.grey[300],
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
            child: Text(
              _selectedItemId == null
                  ? "이용권을 선택해주세요"
                  : "${_currencyFormat.format(_selectedItem['Amount'])} 결제하기",
              style: const TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            ),
          ),
        ),
      ),
    );
  }
}
