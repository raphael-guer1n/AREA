import 'dart:math';
import 'package:flutter/material.dart';
import '../models/area_definitions.dart';

const List<AreaGradient> areaGradients = [
  AreaGradient(from: Color(0xFF002642), to: Color(0xFF0B3C5D)),
  AreaGradient(from: Color(0xFF840032), to: Color(0xFFA33A60)),
  AreaGradient(from: Color(0xFFE59500), to: Color(0xFFF2B344)),
  AreaGradient(from: Color(0xFF5B834D), to: Color(0xFF68915A)),
  AreaGradient(from: Color(0xFF02040F), to: Color(0xFF1B2640)),
];

AreaGradient pickRandomGradient() {
  if (areaGradients.isEmpty) {
    return const AreaGradient(from: Color(0xFF112A46), to: Color(0xFF1C3D63));
  }
  final index = Random().nextInt(areaGradients.length);
  return areaGradients[index];
}

const List<AreaServiceDefinition> areaServiceCatalog = [
  AreaServiceDefinition(
    id: 'timer',
    name: 'Timer',
    actions: [
      AreaActionDefinition(
        id: 'cron_action',
        title: 'cron action',
        label: 'Timer',
        type: 'cron',
        fields: [
          AreaFieldDefinition(
            name: 'delay',
            type: AreaFieldType.number,
            label: 'Delay (seconds)',
            required: true,
            defaultValue: '0',
          ),
        ],
      ),
    ],
  ),
  AreaServiceDefinition(
    id: 'google',
    name: 'Google',
    reactions: [
      AreaReactionDefinition(
        id: 'create_event',
        title: 'create_event',
        label: 'Create Event',
        fields: [
          AreaFieldDefinition(
            name: 'summary',
            type: AreaFieldType.text,
            label: 'Event title',
            required: true,
            defaultValue: '',
          ),
          AreaFieldDefinition(
            name: 'description',
            type: AreaFieldType.text,
            label: 'Event description',
            required: false,
            defaultValue: '',
          ),
          AreaFieldDefinition(
            name: 'start_time',
            type: AreaFieldType.date,
            label: 'Start Time',
            required: true,
            defaultValue: '',
          ),
          AreaFieldDefinition(
            name: 'end_time',
            type: AreaFieldType.date,
            label: 'End Time',
            required: true,
            defaultValue: '',
          ),
          AreaFieldDefinition(
            name: 'calendar',
            type: AreaFieldType.text,
            label: 'Calendar',
            required: true,
            defaultValue: 'primary',
          ),
        ],
      ),
    ],
  ),
  AreaServiceDefinition(
    id: 'github',
    name: 'GitHub',
  ),
];
