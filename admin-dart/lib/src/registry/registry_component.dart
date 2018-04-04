import 'dart:async';
import 'package:angular/angular.dart';
import 'package:angular_components/angular_components.dart';
import 'package:admin_dart/src/info_chip/info_chip.dart';
import 'registry_service.dart';

@Component(
  selector: 'registry',
  templateUrl: 'registry_component.html',
  styleUrls: const ['registry_component.css'],
  directives: const [
    CORE_DIRECTIVES,
    materialDirectives,
    InfoChip
  ]
)
class RegistryComponent extends OnInit {
  final RegistryService registryService;

  RegistryComponent(RegistryService this.registryService);

  bool isAlive = false;
  bool healthCheckInProgress = true;

  @override
  ngOnInit() {
    healthCheck();
  }

  void healthCheck() {
    healthCheckInProgress = true;
    new Future.delayed(const Duration(seconds: 2), () {
      healthCheckInProgress = false;
      isAlive = true;
    });
  }
}