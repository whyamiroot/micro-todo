import 'package:angular/angular.dart';
import 'package:angular_components/angular_components.dart';
import 'package:admin_dart/src/info_chip/info_chip.dart';
import 'service.dart';

@Component(
  selector: 'service',
  templateUrl: 'service_component.html',
  directives: const [CORE_DIRECTIVES, InfoChip, materialDirectives]
)
class ServiceComponent {
  @Input()
  Service service;
}