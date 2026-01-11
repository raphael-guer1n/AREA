import 'package:flutter/material.dart';
import '../models/area_model.dart';
import '../services/area_service.dart';

class AreaProvider extends ChangeNotifier {
  final AreaService _areaService;

  AreaProvider(this._areaService);

  bool _loading = false;
  String? _error;
  String? _message;

  bool get isLoading => _loading;
  String? get error => _error;
  String? get message => _message;

  Future<void> createEvent(CreateEventRequest req) async {
    _loading = true;
    _error = null;
    _message = null;
    notifyListeners();

    try {
      final result = await _areaService.createEvent(req);
      _message = result['message'] ?? 'Event created successfully';
    } catch (e) {
      _error = e.toString();
    }

    _loading = false;
    notifyListeners();
  }
}