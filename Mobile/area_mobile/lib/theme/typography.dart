import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'colors.dart';

/// Typography using Josefin Sans to mirror the Web globals.css.
TextTheme buildTextTheme(AppColorPalette palette) {
  return TextTheme(
    displayLarge: GoogleFonts.josefinSans(
      fontSize: 32,
      fontWeight: FontWeight.w700,
      color: palette.almostBlack,
      letterSpacing: -0.5,
    ),
    titleLarge: GoogleFonts.josefinSans(
      fontSize: 22,
      fontWeight: FontWeight.w700,
      color: palette.deepBlue,
    ),
    titleMedium: GoogleFonts.josefinSans(
      fontSize: 18,
      fontWeight: FontWeight.w600,
      color: palette.deepBlue,
    ),
    bodyLarge: GoogleFonts.josefinSans(
      fontSize: 18,
      color: palette.almostBlack,
      height: 1.35,
    ),
    bodyMedium: GoogleFonts.josefinSans(
      fontSize: 16,
      color: palette.almostBlack,
      height: 1.4,
    ),
    bodySmall: GoogleFonts.josefinSans(
      fontSize: 14,
      color: palette.darkGrey,
    ),
    labelLarge: GoogleFonts.josefinSans(
      fontSize: 16,
      fontWeight: FontWeight.w600,
      color: palette.deepBlue,
      letterSpacing: 0.1,
    ),
  );
}
