import React from 'react';
import { Character, UpdateCharacterRequest } from '../../../types';
import { Modal } from '../../common/Modal';
import { CharacterDetail } from './CharacterDetail';

interface CharacterDetailModalProps {
  isOpen: boolean;
  character: Character;
  onClose: () => void;
  onUpdate?: (id: string, data: UpdateCharacterRequest) => Promise<void>;
  onGenerateReference?: () => void;
}

export const CharacterDetailModal: React.FC<CharacterDetailModalProps> = ({
  isOpen,
  character,
  onClose,
  onUpdate,
  onGenerateReference,
}) => {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title=""
      size="large"
    >
      <CharacterDetail
        character={character}
        onUpdate={onUpdate}
        onGenerateReference={onGenerateReference}
        onClose={onClose}
      />
    </Modal>
  );
};
