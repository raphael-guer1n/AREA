import 'package:flutter/material.dart';
import 'colors.dart';
import 'typography.dart';

/// Global ThemeData matching the AREA Design System (mobile).
/// This ensures color, typography, and component coherence across screens.

final ThemeData areaTheme = ThemeData(
  useMaterial3: true,
  brightness: Brightness.light,
  scaffoldBackgroundColor: AppColors.white,
  primaryColor: AppColors.deepBlue,
  textTheme: AppTextTheme,
  colorScheme: const ColorScheme.light(
    primary: AppColors.deepBlue,
    secondary: AppColors.softSlate,
    background: AppColors.lightGrey,
    surface: AppColors.white,
    onPrimary: AppColors.white,
    onSurface: AppColors.almostBlack,
  ),
  inputDecorationTheme: InputDecorationTheme(
    fillColor: AppColors.lightGrey,
    filled: true,
    border: OutlineInputBorder(
      borderRadius: BorderRadius.circular(8),
      borderSide: const BorderSide(color: AppColors.grey),
    ),
    focusedBorder: OutlineInputBorder(
      borderRadius: BorderRadius.circular(8),
      borderSide: const BorderSide(color: AppColors.deepBlue, width: 1.5),
    ),
    contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
  ),
  elevatedButtonTheme: ElevatedButtonThemeData(
    style: ElevatedButton.styleFrom(
      backgroundColor: AppColors.midBlue,
      foregroundColor: AppColors.white,
      textStyle: const TextStyle(
        fontWeight: FontWeight.w600,
        fontSize: 16,
      ),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
      ),
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
    ),
  ),
  cardTheme: const CardThemeData(
  color: AppColors.white,
  surfaceTintColor: AppColors.white,
  elevation: 2,
  margin: EdgeInsets.symmetric(vertical: 8, horizontal: 12),
  shape: RoundedRectangleBorder(
    borderRadius: BorderRadius.all(Radius.circular(12)),
  ),
),
);