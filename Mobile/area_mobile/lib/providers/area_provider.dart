import 'package:flutter/foundation.dart';

import '../models/area_backend_models.dart';
import '../services/area_service.dart';

class AreaProvider extends ChangeNotifier {
  final AreaService _areaService;

  AreaProvider(this._areaService);

  bool _loading = false;
  String? _error;
  List<AreaDto> _areas = [];

  bool get isLoading => _loading;
  String? get error => _error;
  List<AreaDto> get areas => _areas;

  void clearError() {
    if (_error == null) return;
    _error = null;
    notifyListeners();
  }

  Future<void> loadAreas() async {
    _loading = true;
    _error = null;
    notifyListeners();

    try {
      _areas = await _areaService.getAreas();
    } catch (e) {
      _error = e.toString().replaceAll('Exception: ', '');
      _areas = [];
    } finally {
      _loading = false;
      notifyListeners();
    }
  }

  Future<SaveAreaResult?> saveArea(AreaDto area) async {
    _loading = true;
    _error = null;
    notifyListeners();

    try {
      final res = await _areaService.saveArea(area);
      await loadAreas(); // saveArea returns {} so refresh list
      return res;
    } catch (e) {
      _error = e.toString().replaceAll('Exception: ', '');
      return null;
    } finally {
      _loading = false;
      notifyListeners();
    }
  }

  Future<void> toggleArea(AreaDto area) async {
    _error = null;
    notifyListeners();

    try {
      if (area.id <= 0) throw Exception('Invalid area id');

      if (area.active) {
        await _areaService.deactivateArea(area.id);
      } else {
        await _areaService.activateArea(area.id);
      }

      await loadAreas();
    } catch (e) {
      _error = e.toString().replaceAll('Exception: ', '');
      notifyListeners();
    }
  }

  Future<void> deleteArea(int areaId) async {
    _loading = true;
    _error = null;
    notifyListeners();

    try {
      await _areaService.deleteArea(areaId);
      await loadAreas();
    } catch (e) {
      _error = e.toString().replaceAll('Exception: ', '');
    } finally {
      _loading = false;
      notifyListeners();
    }
  }
}