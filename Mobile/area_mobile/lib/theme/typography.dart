import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'colors.dart';

/// Typography configuration for AREA mobile.
/// Modern sans-serif (Inter) with clarity and readability as core values.

final AppTextTheme = TextTheme(
  displayLarge: GoogleFonts.inter(
    fontSize: 32,
    fontWeight: FontWeight.bold,
    color: AppColors.almostBlack,
  ),
  titleLarge: GoogleFonts.inter(
    fontSize: 22,
    fontWeight: FontWeight.bold,
    color: AppColors.deepBlue,
  ),
  titleMedium: GoogleFonts.inter(
    fontSize: 18,
    fontWeight: FontWeight.w600,
    color: AppColors.deepBlue,
  ),
  bodyLarge: GoogleFonts.inter(
    fontSize: 18,
    color: AppColors.almostBlack,
  ),
  bodyMedium: GoogleFonts.inter(
    fontSize: 16,
    color: AppColors.almostBlack,
  ),
  bodySmall: GoogleFonts.inter(
    fontSize: 14,
    color: AppColors.darkGrey,
  ),
  labelLarge: GoogleFonts.inter(
    fontSize: 16,
    fontWeight: FontWeight.w600,
    color: AppColors.deepBlue,
  ),
);