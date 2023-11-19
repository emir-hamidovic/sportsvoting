import { useCallback, useEffect, useState } from 'react';

export interface Item {
  id: number;
}

export interface SelectionHook {
  handleDeselectOne: (item: string) => void;
  handleSelectOne: (item: string) => void;
  selected: string[];
}

export const useSelection = (items: string[] = []): SelectionHook => {
  const [selected, setSelected] = useState<string[]>([]);

  useEffect(() => {
    setSelected([]);
  }, [items]);

  const handleSelectOne = useCallback((item: string) => {
      setSelected([item]);
  }, []);


  const handleDeselectOne = useCallback((item: string) => {
    setSelected((prevState) => {
      return prevState.filter((_item) => _item !== item);
    });
  }, []);

  return {
    handleDeselectOne,
    handleSelectOne,
    selected
  };
};
