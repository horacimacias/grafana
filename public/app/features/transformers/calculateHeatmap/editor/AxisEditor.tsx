import { useState } from 'react';

import { SelectableValue, StandardEditorProps, VariableOrigin } from '@grafana/data';
import { getTemplateSrv } from '@grafana/runtime';
import { HeatmapCalculationBucketConfig, HeatmapCalculationMode } from '@grafana/schema';
import { HorizontalGroup, RadioButtonGroup, ScaleDistribution } from '@grafana/ui';

import { SuggestionsInput } from '../../suggestionsInput/SuggestionsInput';
import { numberOrVariableValidator } from '../../utils';
import { convertDurationToMilliseconds } from '../utils';

const modeOptions: Array<SelectableValue<HeatmapCalculationMode>> = [
  {
    label: 'Size',
    value: HeatmapCalculationMode.Size,
    description: 'Split the buckets based on size',
  },
  {
    label: 'Count',
    value: HeatmapCalculationMode.Count,
    description: 'Split the buckets based on count',
  },
];

const logModeOptions: Array<SelectableValue<HeatmapCalculationMode>> = [
  {
    label: 'Split',
    value: HeatmapCalculationMode.Size,
    description: 'Split the buckets based on size',
  },
];

export const AxisEditor = ({ value, onChange, item }: StandardEditorProps<HeatmapCalculationBucketConfig>) => {
  const [isInvalid, setInvalid] = useState<boolean>(false);

  const allowInterval = item.settings?.allowInterval ?? false;

  const onValueChange = (bucketValue: string) => {
    let isValid = true;
    if (!allowInterval) {
      isValid = numberOrVariableValidator(bucketValue);
    } else if (bucketValue !== '') {
      let durationMS = convertDurationToMilliseconds(bucketValue);
      if (durationMS === undefined) {
        isValid = false;
      }
    }

    setInvalid(!isValid);
    onChange({
      ...value,
      value: bucketValue,
    });
  };

  const templateSrv = getTemplateSrv();
  const variables = templateSrv.getVariables().map((v) => {
    return { value: v.name, label: v.label || v.name, origin: VariableOrigin.Template };
  });

  return (
    <HorizontalGroup>
      <RadioButtonGroup
        value={value?.mode || HeatmapCalculationMode.Size}
        options={value?.scale?.type === ScaleDistribution.Log ? logModeOptions : modeOptions}
        onChange={(mode) => {
          onChange({
            ...value,
            value: '',
            mode,
          });
        }}
      />
      <SuggestionsInput
        invalid={isInvalid}
        error={'Value needs to be an integer or a variable'}
        value={value?.value ?? ''}
        placeholder="Auto"
        onChange={onValueChange}
        suggestions={variables}
      />
    </HorizontalGroup>
  );
};
