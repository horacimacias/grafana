import { useState } from 'react';

import { Action, DataFrame, VariableSuggestion } from '@grafana/data';

import { Button } from '../Button';
import { Modal } from '../Modal/Modal';

import { ActionEditor } from './ActionEditor';

interface ActionEditorModalContentProps {
  action: Action;
  index: number;
  data: DataFrame[];
  onSave: (index: number, action: Action) => void;
  onCancel: (index: number) => void;
  getSuggestions: () => VariableSuggestion[];
}

export const ActionEditorModalContent = ({
  action,
  index,
  onSave,
  onCancel,
  getSuggestions,
}: ActionEditorModalContentProps) => {
  const [dirtyAction, setDirtyAction] = useState(action);
  return (
    <>
      <ActionEditor
        value={dirtyAction}
        index={index}
        onChange={(index, action) => {
          setDirtyAction(action);
        }}
        suggestions={getSuggestions()}
      />
      <Modal.ButtonRow>
        <Button variant="secondary" onClick={() => onCancel(index)} fill="outline">
          Cancel
        </Button>
        <Button
          onClick={() => {
            onSave(index, dirtyAction);
          }}
        >
          Save
        </Button>
      </Modal.ButtonRow>
    </>
  );
};