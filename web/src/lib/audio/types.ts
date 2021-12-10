export type AudioDevices = {
  cards: Card[];
  sources: CardDevice[];
  sinks: CardDevice[];
  timestamp: number;
};

export type Card = {
  index: number;
  description: string;
  bus: AudioDeviceBus;
  formFactor: AudioDeviceFormFactor;
  sourceIds: Pick<CardDevice, 'index'>;
  sinkIds: Pick<CardDevice, 'index'>;
  profiles: BluetoothAudioDeviceProfile[];
  activeProfile: BluetoothAudioDeviceProfile;
};

export type CardDevice = {
  index: number;
  name: string;
  description: string;
  isDefault: boolean;
  volume: number;
  isMuted: boolean;
  bluetoothProtocol: BluetoothProtocol;
};

export enum AudioDeviceBus {
  PCI = 1,
  Bluetooth = 2,
  USB = 3,
}

export enum AudioDeviceFormFactor {
  Internal = 1,
  Headphones = 2,
  Webcam = 3,
  Headset = 4,
}

export enum BluetoothAudioDeviceProfile {
  HeadsetHeadUnit = 1,
  A2DPSinkSBC = 2,
  A2DPSinkAAC = 3,
  A2DPSinkAptX = 4,
  A2DPSinkAptXHD = 5,
  A2DPSinkLDAC = 6,
  Off = 7,
}

export enum BluetoothProtocol {
  HeadsetHeadUnit = 1,
  A2DPSink = 2,
}

export function bluetoothAudioDeviceProfileToString(profile: BluetoothAudioDeviceProfile): string {
  switch (profile) {
    case BluetoothAudioDeviceProfile.HeadsetHeadUnit:
      return 'HeadsetHeadUnit';
    case BluetoothAudioDeviceProfile.A2DPSinkSBC:
      return 'A2DPSinkSBC';
    case BluetoothAudioDeviceProfile.A2DPSinkAAC:
      return 'A2DPSinkAAC';
    case BluetoothAudioDeviceProfile.A2DPSinkAptX:
      return 'A2DPSinkAptX';
    case BluetoothAudioDeviceProfile.A2DPSinkAptXHD:
      return 'A2DPSinkAptXHD';
    case BluetoothAudioDeviceProfile.A2DPSinkLDAC:
      return 'A2DPSinkLDAC';
    case BluetoothAudioDeviceProfile.Off:
      return 'Off';
  }
}
