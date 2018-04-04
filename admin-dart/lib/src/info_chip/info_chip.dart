import 'package:angular/angular.dart';
import 'package:angular_components/angular_components.dart';

@Component(
  selector: 'info-chip',
  templateUrl: 'info_chip.html',
  directives: const [materialDirectives]
)
class InfoChip {
  @Input()
  String label;
  @Input()
  String value;
}