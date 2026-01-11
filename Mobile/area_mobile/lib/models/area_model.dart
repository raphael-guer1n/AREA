class EventModel {
  String startTime;
  String endTime;
  String summary;
  String description;

  EventModel({
    required this.startTime,
    required this.endTime,
    required this.summary,
    required this.description,
  });

  Map<String, dynamic> toJson() => {
        'startTime': startTime,
        'endTime': endTime,
        'summary': summary,
        'description': description,
      };
}

class CreateEventRequest {
  final int delay;
  final EventModel event;

  CreateEventRequest({
    required this.delay,
    required this.event,
  });

  Map<String, dynamic> toJson() => {
        'delay': delay,
        'event': event.toJson(),
      };
}