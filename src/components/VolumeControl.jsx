function VolumeControl ({
  volume,
  onVolumeChange,
  isSoundEnabled,
  onToggleSound
}) {
  // スライダーが変更されたときに呼ばれる関数
  const handleChange = event => {
    // onVolumeChangeを通じてAppコンポーネントの状態を更新
    // event.target.value は文字列なので数値に変換
    onVolumeChange(parseInt(event.target.value, 10))
  }

  return (
    <div className='volume-control-wrapper'>
      <div className='volume-control'>
        <label htmlFor='volumeSlider'>音量: </label>
        <input
          type='range'
          id='volumeSlider'
          min='0'
          max='100'
          value={volume}
          onChange={handleChange}
          aria-label='音量調節スライダー'
        />
        <span id='volumeValue' className='volume-display'>
          {volume}%
        </span>

        <input
          type='checkbox'
          id='soundToggleCheckbox'
          checked={isSoundEnabled}
          onChange={onToggleSound}
        />
        <label htmlFor='soundToggleCheckbox'>音を有効にする</label>
      </div>
    </div>
  )
}
export default VolumeControl
