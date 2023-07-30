import { useCallback, useEffect, useState } from 'react';

export interface Item {
  id: number;
}

export interface SelectionHook {
  handleDeselectAll: () => void;
  handleDeselectOne: (item: string) => void;
  handleSelectAll: () => void;
  handleSelectOne: (item: string) => void;
  selected: string[];
}

export const useSelection = (items: string[] = []): SelectionHook => {
  const [selected, setSelected] = useState<string[]>([]);

  useEffect(() => {
    setSelected([]);
  }, [items]);

  const handleSelectAll = useCallback(() => {
    setSelected([...items]);
  }, [items]);

  const handleSelectOne = useCallback((item: string) => {
    setSelected((prevState) => [...prevState, item]);
  }, []);

  const handleDeselectAll = useCallback(() => {
    setSelected([]);
  }, []);

  const handleDeselectOne = useCallback((item: string) => {
    setSelected((prevState) => {
      return prevState.filter((_item) => _item !== item);
    });
  }, []);

  return {
    handleDeselectAll,
    handleDeselectOne,
    handleSelectAll,
    handleSelectOne,
    selected
  };
};
